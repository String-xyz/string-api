package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/internal/unit21"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/lib/pq"
)

type Transaction interface {
	Quote(d model.TransactionRequest) (model.ExecutionRequest, error)
	Execute(e model.ExecutionRequest, userId string) (model.TransactionReceipt, error)
	New(repos TransactionRepos) Transaction
}

type TransactionRepos struct {
	Asset       repository.Asset
	Network     repository.Network
	Transaction repository.Transaction
	TxLeg       repository.TxLeg
	User        repository.User
	Instrument  repository.Instrument
	Device      repository.Device
	Location    repository.Location
}

type transactionInstruments struct {
	StringBankId   string
	StringWalletId string
}

type transaction struct {
	repos            TransactionRepos
	instruments      transactionInstruments
	stringUserId     string
	stringDeviceId   string
	stringPlatformId string
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
		return res, common.StringError(err)
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, common.StringError(err)
	}

	estimateUSD, _, err := testTransaction(executor, d, chain, true)
	if err != nil {
		return res, common.StringError(err)
	}
	res.Quote = estimateUSD
	executor.Close()

	// Sign entire payload
	signature, err := common.EVMSign(res)
	if err != nil {
		return res, common.StringError(err)
	}
	res.Signature = signature

	return res, nil
}

func (t transaction) Execute(e model.ExecutionRequest, userId string) (model.TransactionReceipt, error) {
	t.getStringInstrumentsAndUserId()
	res := model.TransactionReceipt{}

	user, err := t.repos.User.GetById(userId)
	if err != nil {
		return res, common.StringError(err)
	}
	if user.ID != userId {
		return res, common.StringError(errors.New("not logged in"))
	}

	// Pull chain info needed for execution from repository
	chain, err := ChainInfo(uint64(e.ChainID), t.repos.Network, t.repos.Asset)
	if err != nil {
		return res, common.StringError(err)
	}

	// Create new Tx in repository, populate it with known info
	db, err := t.repos.Transaction.Create(model.Transaction{Status: "Created", NetworkID: chain.UUID, DeviceID: t.stringDeviceId, PlatformID: t.stringPlatformId})
	if err != nil {
		return res, common.StringError(err)
	}
	updateDB := &model.TransactionUpdates{}
	processingFeeAsset, err := t.populateInitialTxModelData(e, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		fmt.Printf("\nERROR = %+v", err)
		return res, common.StringError(err)
	}

	// Dial the RPC and update model status
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, common.StringError(err)
	}
	status := "RPC Dialed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}

	// Test the Tx and update model status
	estimateUSD, estimateETH, err := testTransaction(executor, e.TransactionRequest, chain, false)
	if err != nil {
		return res, common.StringError(err)
	}
	status = "Tested and Estimated"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}

	// Verify the Quote and update model status
	_, err = verifyQuote(e, estimateUSD)
	if err != nil {
		return res, common.StringError(err)
	}
	status = "Quote Verified"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}

	// Get current balance of primary token
	preBalance, err := executor.GetBalance()
	if err != nil {
		return res, common.StringError(err)
	}
	if preBalance < estimateETH {
		msg := fmt.Sprintf("STRING-API: %s balance is too low to execute %.2f transaction at %.2f", chain.OwlracleName, estimateETH, preBalance)
		MessageStaff(msg)
		return res, common.StringError(errors.New("hot wallet ETH balance too low"))
	}

	// Authorize quoted cost on end-user CC and update model status
	cardAuthorization, err := t.authCard(e.UserAddress, e.CardToken, e.TotalUSD, processingFeeAsset, db.ID, userId)
	if err != nil {
		return res, common.StringError(err)
	}
	status = "Card Authorized"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}

	// Send request to the blockchain and update model status, hash, transaction amount
	txID, value, err := t.initiateTransaction(executor, e, processingFeeAsset, db.ID, userId)
	if err != nil {
		return res, common.StringError(err)
	}
	status = "Transaction Initiated"
	updateDB.Status = &status
	updateDB.TransactionHash = &txID
	txAmount := value.String()
	updateDB.TransactionAmount = &txAmount
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}

	// this Executor will not exist in scope of postProcess
	executor.Close()

	// Send required information to new thread and return TxID to the endpoint
	post := postProcessRequest{
		TxID:               txID,
		Chain:              chain,
		AuthorizationID:    cardAuthorization.AuthID,
		UserAddress:        e.UserAddress,
		CumulativeValue:    value,
		QuotedTotal:        e.TotalUSD,
		TxDBID:             db.ID,
		processingFeeAsset: processingFeeAsset,
		preBalance:         preBalance,
		userId:             userId,
	}
	go t.postProcess(post)

	return model.TransactionReceipt{TxID: txID}, nil
}

func (t *transaction) getStringInstrumentsAndUserId() error {
	// TODO: Look up our instruments and user ID from db
	t.instruments.StringBankId = "13438963-f5e7-47c4-a790-ebca3e3bf915"
	t.instruments.StringWalletId = "ab6a2d66-ad4c-43f4-adf9-c0cd3282492c"
	t.stringUserId = "0e837b73-55cf-43ff-9b1e-0d8258eec978"
	t.stringDeviceId = "073f5a88-9223-4554-a7ce-11d358123a21"
	t.stringPlatformId = "e2724c34-51f6-4eb9-a219-8fb6fb3cbb17"
	return nil
}

func (t transaction) populateInitialTxModelData(e model.ExecutionRequest, m *model.TransactionUpdates) (model.Asset, error) {
	txType := "fiat-to-crypto"
	m.Type = &txType
	// TODO populate db.Tags with key-val pairs for Unit21
	// TODO populate db.DeviceID with info from fingerprint
	// TODO populate db.IPAddress with info from fingerprint
	// TODO populate db.PlatformID with UUID of customer
	// bytes, err := json.Marshal()

	contractParams := pq.StringArray(e.CxParams)
	m.ContractParams = &contractParams
	contractFunc := e.CxFunc + e.CxReturn
	m.ContractFunc = &contractFunc

	asset, err := t.repos.Asset.GetName("USD")
	if err != nil {
		return model.Asset{}, common.StringError(err)
	}
	m.ProcessingFeeAsset = &asset.ID // Checkout processing asset
	return asset, nil
}

func testTransaction(executor Executor, t model.TransactionRequest, chain Chain, useBuffer bool) (model.Quote, float64, error) {
	res := model.Quote{}

	call := ContractCall{
		CxAddr:     t.CxAddr,
		CxFunc:     t.CxFunc,
		CxReturn:   t.CxReturn,
		CxParams:   t.CxParams,
		TxValue:    t.TxValue,
		TxGasLimit: t.TxGasLimit,
	}
	// Estimate value and gas of Tx request
	estimateEVM, err := executor.Estimate(call)
	if err != nil {
		return res, 0, common.StringError(err)
	}

	// Calculate total eth estimate as float64
	gas := new(big.Int)
	gas.SetUint64(estimateEVM.Gas)
	wei := gas.Add(&estimateEVM.Value, gas)
	eth := common.WeiToEther(wei)

	chainID, err := executor.GetChainID()
	if err != nil {
		return res, eth, common.StringError(err)
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
	// Estimate Cost in USD to execute Tx request
	estimateUSD, err := cost.EstimateTransaction(estimationParams, chain)
	if err != nil {
		return res, eth, common.StringError(err)
	}
	res = estimateUSD
	return res, eth, nil
}

func verifyQuote(e model.ExecutionRequest, newEstimate model.Quote) (bool, error) {
	// Null out values which have changed since payload was signed
	dataToValidate := e
	dataToValidate.Signature = ""
	dataToValidate.CardToken = ""
	valid, err := common.ValidateEVMSignature(e.Signature, dataToValidate)
	if err != nil {
		return false, common.StringError(err)
	}
	if !valid {
		return false, common.StringError(errors.New("verifyQuote: invalid signature"))
	}
	if newEstimate.Timestamp-e.Timestamp > 20000 {
		return false, common.StringError(errors.New("verifyQuote: quote expired"))
	}
	if newEstimate.TotalUSD > e.TotalUSD {
		return false, common.StringError(errors.New("verifyQuote: price too volatile"))
	}
	return true, nil
}

func (t transaction) addCardInstrumentIdIfNew(fingerprint string, userID string, last4 string) (string, error) {
	instrument, err := t.repos.Instrument.GetWallet(fingerprint)   // temporarily using get wallet and storing it there
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return "", common.StringError(err)
	} else if err == nil && instrument.UserID != "" {
		return instrument.ID, nil // instrument already exists
	}

	// Create a new instrument
	instrument = model.Instrument{Type: "card", Status: "authorized", Last4: last4, UserID: userID, PublicKey: fingerprint} // No locationID until fingerprint
	instrument, err = t.repos.Instrument.Create(instrument)
	if err != nil {
		return "", common.StringError(err)
	}
	return instrument.ID, nil
}

func (t transaction) addWalletInstrumentIdIfNew(address string) (string, error) {
	instrument, err := t.repos.Instrument.GetWallet(address)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return "", common.StringError(err)
	} else if err == nil && instrument.PublicKey == address {
		return instrument.ID, nil
	}

	// Create a new instrument
	instrument = model.Instrument{Type: "crypto-wallet", Status: "external", Network: "ethereum", PublicKey: address} // No locationID or userID because this wallet was not registered with the user and is some other recipient
	instrument, err = t.repos.Instrument.Create(instrument)
	if err != nil {
		return "", common.StringError(err)
	}
	return instrument.ID, nil
}

func (t transaction) authCard(userWallet string, cardToken string, usd float64, chargeAsset model.Asset, dbID string, userId string) (AuthorizedCharge, error) {
	// auth their card
	auth := AuthorizedCharge{}
	auth, err := AuthorizeCharge(usd, userWallet, cardToken)
	if err != nil {
		return auth, common.StringError(err)
	}

	// Add Checkout Instrument ID to our DB if it's not there already and associate it with the user
	instrumentId, err := t.addCardInstrumentIdIfNew(auth.InstrumentFingerprint, userId, auth.Last4)
	if err != nil {
		return auth, common.StringError(err)
	}

	// Create Origin Tx leg
	usdWei := floatToFixedString(usd, int(chargeAsset.Decimals))
	origin := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetID:      chargeAsset.ID,
		UserID:       userId,
		InstrumentID: instrumentId,
	}
	origin, err = t.repos.TxLeg.Create(origin)
	if err != nil {
		return auth, common.StringError(err)
	}
	txLeg := model.TransactionUpdates{OriginTxLegID: &origin.ID}
	err = t.repos.Transaction.Update(dbID, txLeg)
	if err != nil {
		return auth, common.StringError(err)
	}

	// Send Instrument Data to Unit21
	instrumentModel, err := t.repos.Instrument.GetById(origin.InstrumentID)
	if err != nil {
		return auth, common.StringError(err)
	}

	u21Repo := unit21.InstrumentRepo{
		User:     t.repos.User,
		Device:   t.repos.Device,
		Location: t.repos.Location, // empty until fingerprint integration
	}

	u21Tx := unit21.NewInstrument(u21Repo)
	_, err = u21Tx.Create(instrumentModel)
	if err != nil {
		return auth, common.StringError(err)
	}

	return auth, nil
}

func (t transaction) initiateTransaction(executor Executor, e model.ExecutionRequest, chargeAsset model.Asset, txUUID string, userId string) (string, *big.Int, error) {
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
		return "", nil, common.StringError(err)
	}

	// Create Send Tx leg
	eth := common.WeiToEther(value)
	wei := floatToFixedString(eth, 18)
	usd := floatToFixedString(e.TotalUSD, int(chargeAsset.Decimals))
	send := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetID:      chargeAsset.ID,
		UserID:       userId,
		InstrumentID: t.instruments.StringWalletId,
	}
	send, err = t.repos.TxLeg.Create(send)
	if err != nil {
		return txID, value, common.StringError(err)
	}
	txLeg := model.TransactionUpdates{ResponseTxLegID: &send.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return txID, value, common.StringError(err)
	}

	return txID, value, nil
}

func confirmTx(executor Executor, txID string) (uint64, error) {
	trueGas, err := executor.TxWait(txID)
	if err != nil {
		return 0, common.StringError(err)
	}
	return trueGas, nil
}

func (t transaction) chargeCard(userWallet string, authorizationID string, usd float64, chargeAsset model.Asset, txUUID string, userId string) error {
	_, err := CaptureCharge(usd, userWallet, authorizationID)
	if err != nil {
		return common.StringError(err)
	}

	// Create Receipt Tx leg
	usdWei := floatToFixedString(usd, int(chargeAsset.Decimals))
	receipt := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetID:      chargeAsset.ID,
		UserID:       t.stringUserId,
		InstrumentID: t.instruments.StringBankId,
	}
	receipt, err = t.repos.TxLeg.Create(receipt)
	if err != nil {
		return common.StringError(err)
	}
	txLeg := model.TransactionUpdates{ReceiptTxLegID: &receipt.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

func (t transaction) tenderTransaction(cumulativeValue *big.Int, cumulativeGas uint64, quotedTotal float64, chain Chain, txUUID string, recipientId string, userWalletId string) (float64, error) {
	cost := NewCost(repository.NewCost(nil)) // temporary nil
	trueWei := big.NewInt(0).Add(cumulativeValue, big.NewInt(int64(cumulativeGas)))
	trueEth := common.WeiToEther(trueWei)
	trueUSD, err := cost.LookupUSD(chain.CoingeckoName, trueEth)
	if err != nil {
		return 0, common.StringError(err)
	}
	profit := quotedTotal - trueUSD

	// Create Receive Tx leg
	asset, err := t.repos.Asset.GetName("ETH")
	if err != nil {
		return profit, common.StringError(err)
	}
	wei := floatToFixedString(trueEth, int(asset.Decimals))
	usd := floatToFixedString(quotedTotal, 6)
	send := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetID:      asset.ID,
		UserID:       recipientId,
		InstrumentID: userWalletId,
	}
	send, err = t.repos.TxLeg.Create(send)
	if err != nil {
		return profit, common.StringError(err)
	}
	txLeg := model.TransactionUpdates{DestinationTxLegID: &send.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return profit, common.StringError(err)
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
	preBalance         float64
	userId             string
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

	// confirm the Tx on the EVM, update db status and NetworkFee
	trueGas, err := confirmTx(executor, request.TxID)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	status = "Tx Confirmed"
	updateDB.Status = &status
	networkFee := strconv.FormatUint(trueGas, 10)
	updateDB.NetworkFee = &networkFee // geth uses uint64 for gas
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}

	// Check and see if balance threshold was crossed
	postBalance, err := executor.GetBalance()
	if err != nil {
		// TODO: handle error instead of returning it
	}
	// TODO: store threshold on a per-network basis in the repo
	threshold := 10.0
	if request.preBalance >= threshold && postBalance < threshold {
		msg := fmt.Sprintf("STRING-API: %s balance is < %.2f at %.2f", request.Chain.OwlracleName, threshold, postBalance)
		MessageStaff(msg)
		if err != nil {
			// TODO: handle error instead of returning it
		}
	}

	// compute profit, update db status and processing fees to db
	// TODO: factor request.processingFeeAsset in the event of crypto-to-usd
	recipientWalletId, err := t.addWalletInstrumentIdIfNew(request.UserAddress)
	if err != nil {
		// TODO: handle error instead of returning it
	}
	profit, err := t.tenderTransaction(request.CumulativeValue, trueGas, request.QuotedTotal, request.Chain, request.TxDBID, request.userId, recipientWalletId)
	if err != nil {
		// TODO: Handle error instead of returning it
	}
	fmt.Printf("PROFIT=%+v", profit)
	status = "Profit Tendered"
	updateDB.Status = &status
	stringFee := floatToFixedString(profit, 6)
	processingFee := floatToFixedString(profit, 6) // TODO: set processingFee based on payment method, and location
	updateDB.StringFee = &stringFee                // string fee is always USD with 6 digits
	updateDB.ProcessingFee = &processingFee
	err = t.repos.Transaction.Update(request.TxDBID, updateDB)
	if err != nil {
		// TODO: Handle error instead of returning it
	}

	// charge the users CC
	err = t.chargeCard(request.UserAddress, request.AuthorizationID, request.QuotedTotal, request.processingFeeAsset, request.TxDBID, request.userId)
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
	// Create Transaction data in Unit21
	txModel, err := t.repos.Transaction.GetById(request.TxDBID)
	if err != nil {
		log.Printf("Error getting tx model in Unit21 in Tx Postprocess: %s", err)
		// return res, common.StringError(err)
	}

	u21Repo := unit21.TransactionRepo{
		TxLeg: t.repos.TxLeg,
		User:  t.repos.User,
		Asset: t.repos.Asset,
	}

	u21Tx := unit21.NewTransaction(u21Repo)
	_, err = u21Tx.Create(txModel)
	if err != nil {
		log.Printf("Error updating Unit21 in Tx Postprocess: %s", err)
		// return res, common.StringError(err)
	}
}

func floatToFixedString(value float64, decimals int) string {
	return strconv.FormatUint(uint64(value*(math.Pow10(decimals-1))), 10)
}
