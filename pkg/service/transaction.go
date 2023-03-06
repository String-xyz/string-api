package service

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Transaction interface {
	Quote(d model.TransactionRequest) (model.PrecisionSafeExecutionRequest, error)
	Execute(e model.PrecisionSafeExecutionRequest, userId string, deviceId string, ip string) (model.TransactionReceipt, error)
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
	Contact     repository.Contact
}

type InternalIds struct {
	StringBankId     string `json:"stringBankId" db:"string_bank_id"`
	StringWalletId   string `json:"stringWalletId" db:"string_wallet_id"`
	StringUserId     string `json:"stringUserId" db:"string_user_id"`
	StringPlatformId string `json:"stringPlatformId" db:"string_platform_id"` // temporary
}

type transaction struct {
	repos  repository.Repositories
	redis  store.RedisStore
	ids    InternalIds
	unit21 Unit21
}

func NewTransaction(repos repository.Repositories, redis store.RedisStore, unit21 Unit21) Transaction {
	return &transaction{repos: repos, redis: redis, unit21: unit21}
}

type transactionProcessingData struct {
	userId                        *string
	user                          *model.User
	deviceId                      *string
	ip                            *string
	executor                      *Executor
	processingFeeAsset            *model.Asset
	transactionModel              *model.Transaction
	chain                         *Chain
	executionRequest              *model.ExecutionRequest
	precisionSafeExecutionRequest *model.PrecisionSafeExecutionRequest
	cardAuthorization             *AuthorizedCharge
	cardCapture                   *payments.CapturesResponse
	preBalance                    *float64
	recipientWalletId             *string
	txId                          *string
	cumulativeValue               *big.Int
	trueGas                       *uint64
}

func (t transaction) Quote(d model.TransactionRequest) (model.PrecisionSafeExecutionRequest, error) {
	// TODO: use prefab service to parse d and fill out known params
	res := model.PrecisionSafeExecutionRequest{TransactionRequest: d}
	// chain, err := model.ChainInfo(uint64(d.ChainId))
	chain, err := ChainInfo(uint64(d.ChainId), t.repos.Network, t.repos.Asset)
	if err != nil {
		return res, common.StringError(err)
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, common.StringError(err)
	}

	estimateUSD, _, err := t.testTransaction(executor, d, chain, true)
	if err != nil {
		return res, common.StringError(err)
	}
	res.PrecisionSafeQuote = common.QuoteToPrecise(estimateUSD)
	executor.Close()

	// Sign entire payload
	bytes, err := json.Marshal(res)
	if err != nil {
		return res, common.StringError(err)
	}
	signature, err := common.EVMSign(bytes, true)
	if err != nil {
		return res, common.StringError(err)
	}
	res.Signature = signature

	return res, nil
}

func (t transaction) Execute(e model.PrecisionSafeExecutionRequest, userId string, deviceId string, ip string) (res model.TransactionReceipt, err error) {
	t.getStringInstrumentsAndUserId()
	p := transactionProcessingData{precisionSafeExecutionRequest: &e, executionRequest: &model.ExecutionRequest{}, userId: &userId, deviceId: &deviceId, ip: &ip}

	// Pre-flight transaction setup
	p, err = t.transactionSetup(p)
	if err != nil {
		return res, common.StringError(err)
	}

	// Run safety checks
	p, err = t.safetyCheck(p)
	if err != nil {
		return res, common.StringError(err)
	}

	// Send request to the blockchain and update model status, hash, transaction amount
	p, err = t.initiateTransaction(p)
	if err != nil {
		return res, common.StringError(err)
	}

	// this Executor will not exist in scope of postProcess
	(*p.executor).Close()

	// Send required information to new thread and return txId to the endpoint
	go t.postProcess(p)

	return model.TransactionReceipt{TxId: *p.txId, TxURL: p.chain.Explorer + "/tx/" + *p.txId}, nil
}

func (t transaction) transactionSetup(p transactionProcessingData) (transactionProcessingData, error) {
	// get user object
	user, err := t.repos.User.GetById(*p.userId)
	if err != nil {
		return p, common.StringError(err)
	}
	email, err := t.repos.Contact.GetByUserIdAndType(user.Id, "email")
	if err != nil && errors.Cause(err).Error() != "not found" {
		return p, common.StringError(err)
	}
	user.Email = email.Data
	p.user = &user

	// Pull chain info needed for execution from repository
	chain, err := ChainInfo(p.precisionSafeExecutionRequest.ChainId, t.repos.Network, t.repos.Asset)
	if err != nil {
		return p, common.StringError(err)
	}
	p.chain = &chain

	// Create new Tx in repository, populate it with known info
	transactionModel, err := t.repos.Transaction.Create(model.Transaction{Status: "Created", NetworkId: chain.UUID, DeviceId: *p.deviceId, IPAddress: *p.ip, PlatformId: t.ids.StringPlatformId})
	if err != nil {
		return p, common.StringError(err)
	}
	p.transactionModel = &transactionModel

	updateDB := &model.TransactionUpdates{}
	processingFeeAsset, err := t.populateInitialTxModelData(*p.precisionSafeExecutionRequest, updateDB)
	p.processingFeeAsset = &processingFeeAsset
	if err != nil {
		return p, common.StringError(err)
	}
	err = t.repos.Transaction.Update(transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Send()
		return p, common.StringError(err)
	}

	// Dial the RPC and update model status
	executor := NewExecutor()
	p.executor = &executor
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return p, common.StringError(err)
	}

	err = t.updateTransactionStatus("RPC Dialed", transactionModel.Id)
	if err != nil {
		return p, common.StringError(err)
	}

	return p, err
}

func (t transaction) safetyCheck(p transactionProcessingData) (transactionProcessingData, error) {
	// Test the Tx and update model status
	estimateUSD, estimateETH, err := t.testTransaction(*p.executor, p.precisionSafeExecutionRequest.TransactionRequest, *p.chain, false)
	if err != nil {
		return p, common.StringError(err)
	}
	err = t.updateTransactionStatus("Tested and Estimated", p.transactionModel.Id)
	if err != nil {
		return p, common.StringError(err)
	}

	// Verify the Quote and update model status
	_, err = verifyQuote(*p.precisionSafeExecutionRequest, estimateUSD)
	if err != nil {
		return p, common.StringError(err)
	}
	err = t.updateTransactionStatus("Quote Verified", p.transactionModel.Id)
	if err != nil {
		return p, common.StringError(err)
	}
	*p.executionRequest = common.ExecutionRequestToImprecise(*p.precisionSafeExecutionRequest)

	// Get current balance of primary token
	preBalance, err := (*p.executor).GetBalance()
	p.preBalance = &preBalance
	if err != nil {
		return p, common.StringError(err)
	}
	if preBalance < estimateETH {
		msg := fmt.Sprintf("STRING-API: %s balance is too low to execute %.2f transaction at %.2f", p.chain.OwlracleName, estimateETH, preBalance)
		MessageStaff(msg)
		return p, common.StringError(errors.New("hot wallet ETH balance too low"))
	}

	// Authorize quoted cost on end-user CC and update model status
	p, err = t.authCard(p)
	if err != nil {
		return p, common.StringError(err)
	}

	// Validate Transaction through Real Time Rules engine
	txModel, err := t.repos.Transaction.GetById(p.transactionModel.Id)
	if err != nil {
		log.Err(err).Msg("error getting tx model in unit21 Tx Evalute")
		return p, common.StringError(err)
	}

	evaluation, err := t.unit21.Transaction.Evaluate(txModel)

	if err != nil {
		// If Unit21 Evaluate fails, just log, but otherwise continue with the transaction
		log.Err(err).Msg("Error evaluating transaction in Unit21")
		return p, nil // NOTE: intentionally returning nil here in order to continue the transaction
	}

	if !evaluation {
		err = t.updateTransactionStatus("Failed", p.transactionModel.Id)
		if err != nil {
			return p, common.StringError(err)
		}

		err = t.unit21CreateTransaction(p.transactionModel.Id)
		if err != nil {
			return p, common.StringError(err)
		}

		return p, common.StringError(errors.New("risk: Transaction Failed Unit21 Real Time Rules Evaluation"))
	}

	err = t.updateTransactionStatus("Unit21 Authorized", p.transactionModel.Id)
	if err != nil {
		return p, common.StringError(err)
	}

	return p, nil
}

func (t transaction) initiateTransaction(p transactionProcessingData) (transactionProcessingData, error) {
	call := ContractCall{
		CxAddr:     p.executionRequest.CxAddr,
		CxFunc:     p.executionRequest.CxFunc,
		CxReturn:   p.executionRequest.CxReturn,
		CxParams:   p.executionRequest.CxParams,
		TxValue:    p.executionRequest.TxValue,
		TxGasLimit: p.executionRequest.TxGasLimit,
	}

	txId, value, err := (*p.executor).Initiate(call)
	p.cumulativeValue = value
	if err != nil {
		return p, common.StringError(err)
	}
	p.txId = &txId

	// Create Response Tx leg
	eth := common.WeiToEther(value)
	wei := floatToFixedString(eth, 18)
	usd := floatToFixedString(p.executionRequest.TotalUSD, int(p.processingFeeAsset.Decimals))
	responseLeg := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetId:      p.processingFeeAsset.Id,
		UserId:       *p.userId,
		InstrumentId: t.ids.StringWalletId,
	}
	responseLeg, err = t.repos.TxLeg.Create(responseLeg)
	if err != nil {
		return p, common.StringError(err)
	}
	txLeg := model.TransactionUpdates{ResponseTxLegId: &responseLeg.Id}
	err = t.repos.Transaction.Update(p.transactionModel.Id, txLeg)
	if err != nil {
		return p, common.StringError(err)
	}

	status := "Transaction Initiated"
	txAmount := p.cumulativeValue.String()
	updateDB := &model.TransactionUpdates{Status: &status, TransactionHash: p.txId, TransactionAmount: &txAmount}
	err = t.repos.Transaction.Update(p.transactionModel.Id, updateDB)
	if err != nil {
		return p, common.StringError(err)
	}

	return p, nil
}

func (t transaction) postProcess(p transactionProcessingData) {
	// Reinitialize Executor
	executor := NewExecutor()
	p.executor = &executor
	err := executor.Initialize(p.chain.RPC)
	if err != nil {
		log.Err(err).Msg("Failed to initialized executor in postProcess")
		// TODO: Handle error instead of returning it
	}

	// Update TX Status
	updateDB := model.TransactionUpdates{}
	status := "Post Process RPC Dialed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Post Process RPC Dialed'")
		// TODO: Handle error instead of returning it
	}

	// confirm the Tx on the EVM
	trueGas, err := confirmTx(executor, *p.txId)
	p.trueGas = &trueGas
	if err != nil {
		log.Err(err).Msg("Failed to confirm transaction")
		// TODO: Handle error instead of returning it
	}

	// Update DB status and NetworkFee
	status = "Tx Confirmed"
	updateDB.Status = &status
	networkFee := strconv.FormatUint(trueGas, 10)
	updateDB.NetworkFee = &networkFee // geth uses uint64 for gas
	err = t.repos.Transaction.Update(p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Tx Confirmed'")
		// TODO: Handle error instead of returning it
	}

	// Get new string wallet balance after executing the transaction
	postBalance, err := executor.GetBalance()
	if err != nil {
		log.Err(err).Msg("Failed to get executor balance")
		// TODO: handle error instead of returning it
	}

	// We can close the executor because we aren't using it after this
	executor.Close()

	// If threshold was crossed, notify devs
	// TODO: store threshold on a per-network basis in the repo
	threshold := 10.0
	if *p.preBalance >= threshold && postBalance < threshold {
		msg := fmt.Sprintf("STRING-API: %s balance is < %.2f at %.2f", p.chain.OwlracleName, threshold, postBalance)
		err = MessageStaff(msg)
		if err != nil {
			log.Err(err).Msg("Failed to send staff with low balance threshold message")
			// Not seeing any e
			// TODO: handle error instead of returning it
		}
	}

	// compute profit
	// TODO: factor request.processingFeeAsset in the event of crypto-to-usd
	profit, err := t.tenderTransaction(p)
	if err != nil {
		log.Err(err).Msg("Failed to tender transaction")
		// TODO: Handle error instead of returning it
	}
	stringFee := floatToFixedString(profit, 6)
	processingFee := floatToFixedString(profit, 6) // TODO: set processingFee based on payment method, and location

	// update db status and processing fees to db
	updateDB.StringFee = &stringFee // string fee is always USD with 6 digits
	updateDB.ProcessingFee = &processingFee
	status = "Profit Tendered"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Profit Tendered'")
		// TODO: Handle error instead of returning it
	}

	// charge the users CC
	err = t.chargeCard(p)
	if err != nil {
		log.Err(err).Msg("failed to charge card")
		// TODO: Handle error instead of returning it
	}

	// Update status upon success
	status = "Card Charged"
	updateDB.Status = &status
	// TODO: Figure out how much we paid the CC payment processor and deduct it
	// and use it to populate processing_fee and processing_fee_asset in the table
	err = t.repos.Transaction.Update(p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Card Charged'")
		// TODO: Handle error instead of returning it
	}

	// Transaction complete!  Update status
	status = "Completed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Completed'")
	}

	// Create Transaction data in Unit21
	err = t.unit21CreateTransaction(p.transactionModel.Id)
	if err != nil {
		log.Err(err).Msg("Error creating Unit21 transaction")
	}

	// send email receipt
	err = t.sendEmailReceipt(p)
	if err != nil {
		log.Err(err).Msg("Error sending email receipt to user")
	}
}

func (t transaction) populateInitialTxModelData(e model.PrecisionSafeExecutionRequest, m *model.TransactionUpdates) (model.Asset, error) {
	txType := "fiat-to-crypto"
	m.Type = &txType
	// TODO populate transactionModel.Tags with key-val pairs for Unit21
	// TODO populate transactionModel.DeviceId with info from fingerprint
	// TODO populate transactionModel.IPAddress with info from fingerprint
	// TODO populate transactionModel.PlatformId with UUID of customer
	// bytes, err := json.Marshal()

	contractParams := pq.StringArray(e.CxParams)
	m.ContractParams = &contractParams
	contractFunc := e.CxFunc + e.CxReturn
	m.ContractFunc = &contractFunc

	asset, err := t.repos.Asset.GetByName("USD")
	if err != nil {
		return model.Asset{}, common.StringError(err)
	}
	m.ProcessingFeeAsset = &asset.Id // Checkout processing asset
	return asset, nil
}

func (t transaction) testTransaction(executor Executor, request model.TransactionRequest, chain Chain, useBuffer bool) (model.Quote, float64, error) {
	res := model.Quote{}

	call := ContractCall{
		CxAddr:     request.CxAddr,
		CxFunc:     request.CxFunc,
		CxReturn:   request.CxReturn,
		CxParams:   request.CxParams,
		TxValue:    request.TxValue,
		TxGasLimit: request.TxGasLimit,
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

	chainId, err := executor.GetByChainId()
	if err != nil {
		return res, eth, common.StringError(err)
	}
	cost := NewCost(t.redis)
	estimationParams := EstimationParams{
		ChainId:    chainId,
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

func verifyQuote(e model.PrecisionSafeExecutionRequest, newEstimate model.Quote) (bool, error) {
	// Null out values which have changed since payload was signed
	dataToValidate := e
	dataToValidate.Signature = ""
	dataToValidate.CardToken = ""
	bytesToValidate, err := json.Marshal(dataToValidate)
	if err != nil {
		return false, common.StringError(err)
	}
	valid, err := common.ValidateEVMSignature(e.Signature, bytesToValidate, true)
	if err != nil {
		return false, common.StringError(err)
	}
	if !valid {
		return false, common.StringError(errors.New("verifyQuote: invalid signature"))
	}
	if newEstimate.Timestamp-e.Timestamp > 20 {
		return false, common.StringError(errors.New("verifyQuote: quote expired"))
	}
	quotedTotal, err := strconv.ParseFloat(e.TotalUSD, 64)
	if err != nil {
		return false, common.StringError(err)
	}
	if newEstimate.TotalUSD > quotedTotal {
		return false, common.StringError(errors.New("verifyQuote: price too volatile"))
	}
	return true, nil
}

func (t transaction) addCardInstrumentIdIfNew(p transactionProcessingData) (string, error) {
	instrument, err := t.repos.Instrument.GetCardByFingerprint(p.cardAuthorization.CheckoutFingerprint)
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return "", common.StringError(err)
	} else if err == nil && instrument.UserId != "" {
		go t.unit21.Instrument.Update(instrument) // if instrument already exists, update it anyways
		return instrument.Id, nil                 // return if instrument already exists
	}

	// We should gather type from the payment processor
	instrument_type := "Debit Card"
	if p.cardAuthorization.CardType == "CREDIT" {
		instrument_type = "Credit Card"
	}
	// Create a new instrument
	instrument = model.Instrument{ // No locationId until fingerprint
		Type:      instrument_type,
		Status:    "created",
		Last4:     p.cardAuthorization.Last4,
		UserId:    *p.userId,
		PublicKey: p.cardAuthorization.CheckoutFingerprint,
	}
	instrument, err = t.repos.Instrument.Create(instrument)
	if err != nil {
		return "", common.StringError(err)
	}

	go t.unit21.Instrument.Create(instrument)

	return instrument.Id, nil
}

func (t transaction) addWalletInstrumentIdIfNew(address string, id string) (string, error) {
	instrument, err := t.repos.Instrument.GetWalletByAddr(address)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return "", common.StringError(err)
	} else if err == nil && instrument.PublicKey == address {
		go t.unit21.Instrument.Update(instrument) // if instrument already exists, update it anyways
		return instrument.Id, nil                 // return if instrument already exists
	}

	// Create a new instrument
	instrument = model.Instrument{Type: "Crypto Wallet", Status: "external", Network: "ethereum", PublicKey: address, UserId: id} // No locationId or userId because this wallet was not registered with the user and is some other recipient
	instrument, err = t.repos.Instrument.Create(instrument)
	if err != nil {
		return "", common.StringError(err)
	}

	go t.unit21.Instrument.Create(instrument)

	return instrument.Id, nil
}

func (t transaction) authCard(p transactionProcessingData) (transactionProcessingData, error) {
	// auth their card
	p, err := AuthorizeCharge(p)
	if err != nil {
		return p, common.StringError(err)
	}

	// Add Checkout Instrument ID to our DB if it's not there already and associate it with the user
	instrumentId, err := t.addCardInstrumentIdIfNew(p)
	if err != nil {
		return p, common.StringError(err)
	}

	// Create Origin Tx leg
	usdWei := floatToFixedString(p.executionRequest.TotalUSD, int(p.processingFeeAsset.Decimals))
	origin := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetId:      p.processingFeeAsset.Id,
		UserId:       *p.userId,
		InstrumentId: instrumentId,
	}
	origin, err = t.repos.TxLeg.Create(origin)
	if err != nil {
		return p, common.StringError(err)
	}
	txLegUpdates := model.TransactionUpdates{OriginTxLegId: &origin.Id}
	err = t.repos.Transaction.Update(p.transactionModel.Id, txLegUpdates)
	if err != nil {
		return p, common.StringError(err)
	}

	err = t.updateTransactionStatus("Card "+p.cardAuthorization.Status, p.transactionModel.Id)
	if err != nil {
		return p, common.StringError(err)
	}

	recipientWalletId, err := t.addWalletInstrumentIdIfNew(p.executionRequest.UserAddress, *p.userId)
	p.recipientWalletId = &recipientWalletId
	if err != nil {
		return p, common.StringError(err)
	}

	// TODO: Determine the output of the transaction (destination leg) with Tracers
	destinationLeg := model.TxLeg{
		Timestamp:    time.Now(),         // Required by the db. Should be updated when the tx occurs
		Amount:       "0",                // Required by Unit21. The amount of the asset received by the user
		Value:        "0",                // Default to '0'. The value of the asset received by the user
		AssetId:      p.chain.GasTokenId, // Required by the db. the asset received by the user
		UserId:       *p.userId,          // the user who received the asset
		InstrumentId: recipientWalletId,  // Required by the db. the instrument which received the asset (wallet usually)
	}

	destinationLeg, err = t.repos.TxLeg.Create(destinationLeg)
	if err != nil {
		return p, common.StringError(err)
	}

	txLegUpdates = model.TransactionUpdates{DestinationTxLegId: &destinationLeg.Id}

	err = t.repos.Transaction.Update(p.transactionModel.Id, txLegUpdates)
	if err != nil {
		return p, common.StringError(err)
	}

	if !p.cardAuthorization.Approved {
		err := t.unit21CreateTransaction(p.transactionModel.Id)
		if err != nil {
			return p, common.StringError(err)
		}

		return p, common.StringError(errors.New("payment: Authorization Declined by Checkout"))
	}

	return p, nil
}

func confirmTx(executor Executor, txId string) (uint64, error) {
	trueGas, err := executor.TxWait(txId)
	if err != nil {
		return 0, common.StringError(err)
	}
	return trueGas, nil
}

// TODO: rewrite this transaction to reference the asset(s) received by the user, not what we paid
func (t transaction) tenderTransaction(p transactionProcessingData) (float64, error) {
	cost := NewCost(t.redis)
	trueWei := big.NewInt(0).Add(p.cumulativeValue, big.NewInt(int64(*p.trueGas)))
	trueEth := common.WeiToEther(trueWei)
	trueUSD, err := cost.LookupUSD(p.chain.CoingeckoName, trueEth)
	if err != nil {
		return 0, common.StringError(err)
	}
	profit := p.executionRequest.Quote.TotalUSD - trueUSD

	// Create Receive Tx leg
	asset, err := t.repos.Asset.GetById(p.chain.GasTokenId)
	if err != nil {
		return profit, common.StringError(err)
	}
	wei := floatToFixedString(trueEth, int(asset.Decimals))
	usd := floatToFixedString(p.executionRequest.Quote.TotalUSD, 6)

	txModel, err := t.repos.Transaction.GetById(p.transactionModel.Id)
	if err != nil {
		return profit, common.StringError(err)
	}

	now := time.Now()
	destinationLeg := model.TxLegUpdates{
		Timestamp:    &now,                // updated based on *when the transaction occured* not time.Now()
		Amount:       &wei,                // Should be the amount of the asset received by the user
		Value:        &usd,                // The value of the asset received by the user
		AssetId:      &asset.Id,           // the asset received by the user
		UserId:       p.userId,            // the user who received the asset
		InstrumentId: p.recipientWalletId, // the instrument which received the asset (wallet usually)
	}

	// We now update the destination leg instead of creating it
	err = t.repos.TxLeg.Update(txModel.DestinationTxLegId, destinationLeg)
	if err != nil {
		return profit, common.StringError(err)
	}

	return profit, nil
}

func (t transaction) chargeCard(p transactionProcessingData) error {
	p, err := CaptureCharge(p)
	if err != nil {
		return common.StringError(err)
	}

	// Create Receipt Tx leg
	usdWei := floatToFixedString(p.executionRequest.Quote.TotalUSD, int(p.processingFeeAsset.Decimals))
	receiptLeg := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetId:      p.processingFeeAsset.Id,
		UserId:       t.ids.StringUserId,
		InstrumentId: t.ids.StringBankId,
	}
	receiptLeg, err = t.repos.TxLeg.Create(receiptLeg)
	if err != nil {
		return common.StringError(err)
	}
	txLeg := model.TransactionUpdates{ReceiptTxLegId: &receiptLeg.Id, PaymentCode: &p.cardCapture.Accepted.ActionID}
	err = t.repos.Transaction.Update(p.transactionModel.Id, txLeg)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

func (t transaction) sendEmailReceipt(p transactionProcessingData) error {
	user, err := t.repos.User.GetById(*p.userId)
	if err != nil {
		log.Err(err).Msg("Error getting user from repo")
		return common.StringError(err)
	}
	contact, err := t.repos.Contact.GetByUserId(user.Id)
	if err != nil {
		log.Err(err).Msg("Error getting user contact from repo")
		return common.StringError(err)
	}
	name := user.FirstName // + " " + user.MiddleName + " " + user.LastName
	if name == "" {
		name = "User"
	}
	receiptParams := common.ReceiptGenerationParams{
		ReceiptType:       "NFT Purchase", // TODO: retrieve dynamically
		CustomerName:      name,
		StringPaymentId:   p.transactionModel.Id,
		PaymentDescriptor: "String Digital Asset", // TODO: retrieve dynamically
		TransactionDate:   time.Now().Format(time.RFC1123),
	}
	receiptBody := [][2]string{
		{"Transaction ID", "<a href='" + p.chain.Explorer + "/tx/" + *p.txId + "'>" + *p.txId + "</a>"},
		{"Destination Wallet", "<a href='" + p.chain.Explorer + "/address/" + p.executionRequest.UserAddress + "'>" + p.executionRequest.UserAddress + "</a>"},
		{"Payment Descriptor", receiptParams.PaymentDescriptor},
		{"Payment Method", p.cardAuthorization.Issuer + " " + p.cardAuthorization.Last4},
		{"Platform", "String Demo"},            // TODO: retrieve dynamically
		{"Item Ordered", "String Fighter NFT"}, // TODO: retrieve dynamically
		{"Token ID", "1234"},                   // TODO: retrieve dynamically, maybe after building token transfer detection
		{"Subtotal", common.FloatToUSDString(p.executionRequest.Quote.BaseUSD + p.executionRequest.Quote.TokenUSD)},
		{"Network Fee:", common.FloatToUSDString(p.executionRequest.Quote.GasUSD)},
		{"Processing Fee", common.FloatToUSDString(p.executionRequest.Quote.ServiceUSD)},
		{"Total Charge", common.FloatToUSDString(p.executionRequest.Quote.TotalUSD)},
	}
	err = common.EmailReceipt(contact.Data, receiptParams, receiptBody)
	if err != nil {
		log.Err(err).Msg("Error sending email receipt to user")
		return common.StringError(err)
	}
	return nil
}

func floatToFixedString(value float64, decimals int) string {
	return strconv.FormatUint(uint64(value*(math.Pow10(decimals))), 10)
}

func (t transaction) unit21CreateTransaction(transactionId string) (err error) {
	txModel, err := t.repos.Transaction.GetById(transactionId)
	if err != nil {
		log.Err(err).Msg("Error getting tx model in Unit21 in Tx Postprocess")
		return common.StringError(err)
	}

	_, err = t.unit21.Transaction.Create(txModel)
	if err != nil {
		log.Err(err).Msg("Error updating unit21 in Tx Postprocess")
		return common.StringError(err)
	}

	return nil
}

func (t transaction) updateTransactionStatus(status string, transactionId string) (err error) {
	updateDB := &model.TransactionUpdates{Status: &status}
	err = t.repos.Transaction.Update(transactionId, updateDB)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

func (t *transaction) getStringInstrumentsAndUserId() {
	t.ids = GetStringIdsFromEnv()
}
