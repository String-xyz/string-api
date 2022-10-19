package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

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
	chain, err := common.ChainInfo(uint64(d.ChainID), t.repos.Network, t.repos.Asset)
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
	chain, err := common.ChainInfo(uint64(e.ChainID), t.repos.Network, t.repos.Asset)
	if err != nil {
		return res, err
	}
	fmt.Printf("\nGot Chain Info")

	// Create new TX in repository, populate it with known info
	db, err := t.repos.Transaction.Create(model.Transaction{Status: "Created", NetworkID: chain.UUID})
	if err != nil {
		return res, err
	}
	updateDB := &model.TransactionUpdates{}
	fmt.Printf("\nCreated TX in db")
	processingFeeAsset, err := t.populateInitialTxModelData(e, updateDB)
	if err != nil {
		return res, err
	}
	fmt.Printf("\nGot Initial TX model data")
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		fmt.Printf("\nERROR = %+v", err)
		return res, err
	}
	fmt.Printf("\nUpdate TX table")

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
	authorizationID, err := authCard(e.UserAddress, e.CardToken, e.TotalUSD)
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
	txID, value, err := initiateTransaction(executor, e)
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

func testTransaction(executor Executor, t model.TransactionRequest, chain common.Chain, useBuffer bool) (model.Quote, error) {
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

func authCard(userWallet string, cardToken string, usd float64) (string, error) {
	// auth their card
	auth, err := AuthorizeCharge(usd, userWallet, cardToken)
	return auth, err
}

func initiateTransaction(executor Executor, e model.ExecutionRequest) (string, *big.Int, error) {
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
	return txID, value, nil
}

func confirmTX(executor Executor, txID string) (uint64, error) {
	trueGas, err := executor.TxWait(txID)
	if err != nil {
		return 0, err
	}
	return trueGas, nil
}

func chargeCard(userWallet string, authorizationID string, usd float64) error {
	_, err := CaptureCharge(usd, userWallet, authorizationID)
	return err
}

func tenderTransaction(cumulativeValue *big.Int, cumulativeGas uint64, quotedTotal float64, chain common.Chain) (float64, error) {
	cost := NewCost(repository.NewCost(nil)) // temporary nil
	trueWei := big.NewInt(0).Add(cumulativeValue, big.NewInt(int64(cumulativeGas)))
	trueEth := common.WeiToEther(trueWei)
	trueUSD, err := cost.LookupUSD(chain.CoingeckoName, trueEth)
	if err != nil {
		return 0, err
	}
	profit := quotedTotal - trueUSD
	return profit, nil
}

type postProcessRequest struct {
	TxID               string
	Chain              common.Chain
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
	profit, err := tenderTransaction(request.CumulativeValue, trueGas, request.QuotedTotal, request.Chain)
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
	err = chargeCard(request.UserAddress, request.AuthorizationID, request.QuotedTotal)
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
