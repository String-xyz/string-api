package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/jmoiron/sqlx/types"
)

type Transaction interface {
	Quote(d model.TransactionRequest) (model.ExecutionRequest, error)
	Execute(e model.ExecutionRequest) (model.TransactionReceipt, error)
	New(repos TransactionRepos) Transaction
}

type TransactionRepos struct {
	Asset       repository.Asset
	Network     repository.Network
	Transaction repository.Transaction
	TxLeg       repository.TxLeg
	User        repository.User
}

type transaction struct {
	repos TransactionRepos
}

func (t transaction) New(repos TransactionRepos) Transaction {
	return &transaction{repos: repos}
}

func NewTransaction(repos TransactionRepos) Transaction {
	return &transaction{repos: repos}
}

func (t transaction) Quote(d model.TransactionRequest) (model.ExecutionRequest, error) {
	// TODO: use prefab service to parse d and fill out known params
	res := model.ExecutionRequest{TransactionRequest: d}
	// chain, err := model.ChainInfo(uint64(d.ChainID))
	chain, err := ChainInfo(uint64(d.ChainID), t.repos.Network, t.repos.Asset)
	if err != nil {
		return res, err
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, err
	}

	estimateUSD, err := testTransaction(executor, d, chain, true)
	if err != nil {
		return res, err
	}
	res.Quote = estimateUSD
	executor.Close()

	// Sign entire payload
	signature, err := common.EVMSign(res)
	if err != nil {
		return res, err
	}
	res.Signature = signature

	return res, nil
}

func (t transaction) Execute(e model.ExecutionRequest) (model.TransactionReceipt, error) {
	res := model.TransactionReceipt{}

	// Pull chain info needed for execution from repository
	chain, err := ChainInfo(uint64(e.ChainID), t.repos.Network, t.repos.Asset)
	if err != nil {
		return res, err
	}

	// Create new TX in repository, populate it with known info
	db, err := t.repos.Transaction.Create(model.Transaction{Status: "Created", NetworkID: chain.UUID})
	if err != nil {
		return res, err
	}
	updateDB := &model.TransactionUpdates{}
	processingFeeAsset, err := t.populateInitialTxModelData(e, updateDB)
	if err != nil {
		return res, err
	}
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		fmt.Printf("\nERROR = %+v", err)
		return res, err
	}

	// Dial the RPC and update model status
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, err
	}
	status := "RPC Dialed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, err
	}

	// Test the TX and update model status
	estimateUSD, err := testTransaction(executor, e.TransactionRequest, chain, false)
	if err != nil {
		return res, err
	}
	status = "Tested and Estimated"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, err
	}

	// Verify the Quote and update model status
	_, err = verifyQuote(e, estimateUSD)
	if err != nil {
		return res, err
	}
	status = "Quote Verified"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, err
	}

	// Authorize quoted cost on end-user CC and update model status
	authorizationID, err := t.authCard(e.UserAddress, e.CardToken, e.TotalUSD, processingFeeAsset, db.ID)
	if err != nil {
		return res, err
	}
	status = "Card Authorized"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, err
	}

	// Send request to the blockchain and update model status, hash, transaction amount
	txID, value, err := t.initiateTransaction(executor, e, processingFeeAsset, db.ID)
	if err != nil {
		return res, err
	}
	status = "Transaction Initiated"
	updateDB.Status = &status
	updateDB.TransactionHash = &txID
	txAmount := value.String()
	updateDB.TransactionAmount = &txAmount
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, err
	}

	// this Executor will not exist in scope of postProcess
	executor.Close()

	// Send required information to new thread and return TXID to the endpoint
	post := postProcessRequest{
		TxID:               txID,
		Chain:              chain,
		AuthorizationID:    authorizationID,
		UserAddress:        e.UserAddress,
		CumulativeValue:    value,
		QuotedTotal:        e.TotalUSD,
		TxDBID:             db.ID,
		processingFeeAsset: processingFeeAsset,
	}
	go t.postProcess(post)
	return model.TransactionReceipt{TxID: txID}, nil
}

func (t transaction) populateInitialTxModelData(e model.ExecutionRequest, m *model.TransactionUpdates) (model.Asset, error) {
	txType := "fiat-to-crypto"
	m.Type = &txType
	// TODO populate db.Tags with key-val pairs for Unit21
	// TODO populate db.DeviceID with info from fingerprint
	// TODO populate db.IPAddress with info from fingerprint
	// TODO populate db.PlatformID with UUID of customer
	bytes, err := json.Marshal(e.CxParams)
	if err != nil {
		return model.Asset{}, err
	}
	contractParams := types.JSONText(bytes)
	m.ContractParams = &contractParams
	contractFunc := e.CxFunc + e.CxReturn
	m.ContractFunc = &contractFunc

	asset, err := t.repos.Asset.GetName("USD")
	if err != nil {
		return model.Asset{}, err
	}
	m.ProcessingFeeAsset = &asset.ID // Checkout processing asset
	return asset, nil
}

func testTransaction(executor Executor, t model.TransactionRequest, chain Chain, useBuffer bool) (model.Quote, error) {
	res := model.Quote{}

	call := ContractCall{
		CxAddr:     t.CxAddr,
		CxFunc:     t.CxFunc,
		CxReturn:   t.CxReturn,
		CxParams:   t.CxParams,
		TxValue:    t.TxValue,
		TxGasLimit: t.TxGasLimit,
	}
	// Estimate value and gas of TX request
	estimateEVM, err := executor.Estimate(call)
	if err != nil {
		return res, err
	}

	chainID, err := executor.GetChainID()
	if err != nil {
		return res, err
	}
	cost := NewCost(repository.NewCost(nil))
	estimationParams := EstimationParams{
		ChainID:    chainID,
		CostETH:    estimateEVM.Value,
		UseBuffer:  useBuffer,
		GasUsedWei: estimateEVM.Gas,
		CostToken:  *big.NewInt(0),
		TokenName:  "",
	}
	// Estimate Cost in USD to execute TX request
	estimateUSD, err := cost.EstimateTransaction(estimationParams, chain)
	if err != nil {
		return res, err
	}
	res = estimateUSD
	return res, nil
}

func verifyQuote(e model.ExecutionRequest, newEstimate model.Quote) (bool, error) {
	// Null out values which have changed since payload was signed
	dataToValidate := e
	dataToValidate.Signature = ""
	dataToValidate.CardToken = ""
	valid, err := common.ValidateEVMSignature(e.Signature, dataToValidate)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("verifyQuote: invalid signature")
	}
	if newEstimate.Timestamp-e.Timestamp > 20000 {
		return false, errors.New("verifyQuote: quote expired")
	}
	if newEstimate.TotalUSD > e.TotalUSD {
		return false, errors.New("verifyQuote: price too volatile")
	}
	return true, nil
}

func (t transaction) authCard(userWallet string, cardToken string, usd float64, chargeAsset model.Asset, dbID string) (string, error) {
	// auth their card
	auth, err := AuthorizeCharge(usd, userWallet, cardToken)
	if err != nil {
		return "", err
	}

	// Create Origin TX leg
	usdWei := floatToFixedString(usd, int(chargeAsset.Decimals))
	origin := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetID:      chargeAsset.ID,
		UserID:       "0e837b73-55cf-43ff-9b1e-0d8258eec978", // TODO: Get dynamically
		InstrumentID: "13438963-f5e7-47c4-a790-ebca3e3bf915", // TODO: Get dynamically
	}
	origin, err = t.repos.TxLeg.Create(origin)
	if err != nil {
		return auth, err
	}
	txLeg := model.TransactionUpdates{OriginTXLegID: &origin.ID}
	err = t.repos.Transaction.Update(dbID, txLeg)
	if err != nil {
		return auth, err
	}
	return auth, err
}

func (t transaction) initiateTransaction(executor Executor, e model.ExecutionRequest, chargeAsset model.Asset, txUUID string) (string, *big.Int, error) {
	call := ContractCall{
		CxAddr:     e.CxAddr,
		CxFunc:     e.CxFunc,
		CxReturn:   e.CxReturn,
		CxParams:   e.CxParams,
		TxValue:    e.TxValue,
		TxGasLimit: e.TxGasLimit,
	}
	txID, value, err := executor.Initiate(call)
	if err != nil {
		return "", nil, err
	}

	// Create Send TX leg
	eth := common.WeiToEther(value)
	wei := floatToFixedString(eth, 18)
	usd := floatToFixedString(e.TotalUSD, int(chargeAsset.Decimals))
	send := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetID:      chargeAsset.ID,
		UserID:       "0e837b73-55cf-43ff-9b1e-0d8258eec978", // TODO: Get dynamically
		InstrumentID: "ab6a2d66-ad4c-43f4-adf9-c0cd3282492c", // TODO: Get dynamically
	}
	send, err = t.repos.TxLeg.Create(send)
	if err != nil {
		return txID, value, err
	}
	txLeg := model.TransactionUpdates{ResponseTXLegID: &send.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return txID, value, err
	}

	return txID, value, nil
}

func confirmTX(executor Executor, txID string) (uint64, error) {
	trueGas, err := executor.TxWait(txID)
	if err != nil {
		return 0, err
	}
	return trueGas, nil
}

func (t transaction) chargeCard(userWallet string, authorizationID string, usd float64, chargeAsset model.Asset, txUUID string) error {
	_, err := CaptureCharge(usd, userWallet, authorizationID)
	if err != nil {
		return err
	}

	// Create Receipt TX leg
	usdWei := floatToFixedString(usd, int(chargeAsset.Decimals))
	receipt := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetID:      chargeAsset.ID,
		UserID:       "0e837b73-55cf-43ff-9b1e-0d8258eec978", // TODO: Get dynamically
		InstrumentID: "13438963-f5e7-47c4-a790-ebca3e3bf915", // TODO: Get dynamically
	}
	receipt, err = t.repos.TxLeg.Create(receipt)
	if err != nil {
		return err
	}
	txLeg := model.TransactionUpdates{ReceiptTXLegID: &receipt.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return err
	}

	return nil
}

func (t transaction) tenderTransaction(cumulativeValue *big.Int, cumulativeGas uint64, quotedTotal float64, chain Chain, txUUID string) (float64, error) {
	cost := NewCost(repository.NewCost(nil)) // temporary nil
	trueWei := big.NewInt(0).Add(cumulativeValue, big.NewInt(int64(cumulativeGas)))
	trueEth := common.WeiToEther(trueWei)
	trueUSD, err := cost.LookupUSD(chain.CoingeckoName, trueEth)
	if err != nil {
		return 0, err
	}
	profit := quotedTotal - trueUSD

	// Create Receive TX leg
	asset, err := t.repos.Asset.GetName("ETH")
	if err != nil {
		return profit, err
	}
	wei := floatToFixedString(trueEth, int(asset.Decimals))
	usd := floatToFixedString(quotedTotal, 6)
	send := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetID:      asset.ID,
		UserID:       "0e837b73-55cf-43ff-9b1e-0d8258eec978", // TODO: Get dynamically
		InstrumentID: "ab6a2d66-ad4c-43f4-adf9-c0cd3282492c", // TODO: Get dynamically
	}
	send, err = t.repos.TxLeg.Create(send)
	if err != nil {
		return profit, err
	}
	txLeg := model.TransactionUpdates{DestinationTXLegID: &send.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return profit, err
	}

	return profit, nil
}

type postProcessRequest struct {
	TxID               string
	Chain              Chain
	AuthorizationID    string
	UserAddress        string
	CumulativeGas      uint64
	CumulativeValue    *big.Int
	QuotedTotal        float64
	TxDBID             string
	processingFeeAsset model.Asset
}

func (t transaction) postProcess(request postProcessRequest) {
	executor := NewExecutor()
	err := executor.Initialize(request.Chain.RPC)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	updateDB := model.TransactionUpdates{}
	status := "Post Process RPC Dialed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}

	// confirm the TX on the EVM, update db status and NetworkFee
	trueGas, err := confirmTX(executor, request.TxID)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	status = "TX Confirmed"
	updateDB.Status = &status
	networkFee := strconv.FormatUint(trueGas, 10)
	updateDB.NetworkFee = &networkFee // geth uses uint64 for gas
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}

	// compute profit, update db status and processing fees to db
	// TODO: factor request.processingFeeAsset in the event of crypto-to-usd
	profit, err := t.tenderTransaction(request.CumulativeValue, trueGas, request.QuotedTotal, request.Chain, request.TxDBID)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	fmt.Printf("PROFIT=%+v", profit)
	status = "Profit Tendered"
	updateDB.Status = &status
	stringFee := floatToFixedString(profit, 6)
	updateDB.StringFee = &stringFee // string fee is always USD with 6 digits
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}

	// charge the users CC
	err = t.chargeCard(request.UserAddress, request.AuthorizationID, request.QuotedTotal, request.processingFeeAsset, request.TxDBID)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	status = "Card Charged"
	updateDB.Status = &status
	// TODO: Figure out how much we paid the CC payment processor and deduct it
	// and use it to populate processing_fee and processing_fee_asset in the table
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}

	status = "Completed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	executor.Close()
}

func floatToFixedString(value float64, decimals int) string {
	return strconv.FormatUint(uint64(value*(math.Pow10(decimals-1))), 10)
}
