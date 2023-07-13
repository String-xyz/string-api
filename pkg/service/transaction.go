package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/lmittmann/w3"

	"github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/internal/emailer"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/String-xyz/string-api/pkg/model"
	repository "github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Quote(ctx context.Context, d model.TransactionRequest, platformId string) (res model.Quote, err error)
	Execute(ctx context.Context, e model.ExecutionRequest, userId string, deviceId string, platformId string, ip string) (res model.TransactionReceipt, err error)
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
	Contract    repository.Contract
}

type InternalIds struct {
	StringBankId     string `json:"stringBankId" db:"string_bank_id"`
	StringWalletId   string `json:"stringWalletId" db:"string_wallet_id"`
	StringUserId     string `json:"stringUserId" db:"string_user_id"`
	StringPlatformId string `json:"stringPlatformId" db:"string_platform_id"` // temporary
}

type transaction struct {
	repos  repository.Repositories
	redis  database.RedisStore
	ids    InternalIds
	unit21 Unit21
}

func NewTransaction(repos repository.Repositories, redis database.RedisStore, unit21 Unit21) Transaction {
	return &transaction{repos: repos, redis: redis, unit21: unit21}
}

type transactionProcessingData struct {
	userId             *string
	user               *model.User
	deviceId           *string
	ip                 *string
	platformId         *string
	executor           *Executor
	processingFeeAsset *model.Asset
	transactionModel   *model.Transaction
	chain              *Chain
	executionRequest   *model.ExecutionRequest
	floatEstimate      *model.Estimate[float64]
	cardAuthorization  *AuthorizedCharge
	PaymentStatus      checkout.PaymentStatus
	PaymentId          string
	recipientWalletId  *string
	txIds              []string
	forwardTxIds       []string
	cumulativeValue    *big.Int
	trueGas            uint64
	tokenIds           string
	tokenQuantities    string
}

func (t transaction) Quote(ctx context.Context, d model.TransactionRequest, platformId string) (res model.Quote, err error) {
	_, finish := Span(ctx, "service.transaction.Quote", SpanTag{"platformId": platformId})
	defer finish()

	// TODO: use prefab service to parse d and fill out known params
	res.TransactionRequest = d
	chain, err := ChainInfo(ctx, uint64(d.ChainId), t.repos.Network, t.repos.Asset)
	if err != nil {
		return res, libcommon.StringError(err)
	}

	allowed, err := t.isContractAllowed(ctx, platformId, chain.UUID, d)
	if err != nil {
		return res, libcommon.StringError(err)
	}
	if !allowed {
		return res, libcommon.StringError(serror.CONTRACT_NOT_ALLOWED)
	}

	executor := NewExecutor()
	err = executor.Initialize(chain)
	if err != nil {
		return res, libcommon.StringError(err)
	}

	estimateUSD, _, _, err := t.testTransaction(executor, d, chain, true, true)
	if err != nil {
		return res, libcommon.StringError(err)
	}
	res.Estimate = common.EstimateToPrecise(estimateUSD)
	executor.Close()

	// Sign entire payload
	bytes, err := json.Marshal(res)
	if err != nil {
		return res, libcommon.StringError(err)
	}
	signature, err := common.EVMSign(bytes, true)
	if err != nil {
		return res, libcommon.StringError(err)
	}
	res.Signature = signature

	return res, nil
}

func (t transaction) Execute(ctx context.Context, e model.ExecutionRequest, userId string, deviceId string, platformId string, ip string) (res model.TransactionReceipt, err error) {
	_, finish := Span(ctx, "service.transaction.Execute", SpanTag{"platformId": platformId})
	defer finish()

	t.getStringInstrumentsAndUserId()
	p := transactionProcessingData{executionRequest: &e, userId: &userId, deviceId: &deviceId, ip: &ip, platformId: &platformId}

	// Pre-flight transaction setup
	p, err = t.transactionSetup(ctx, p)
	if err != nil {
		return res, libcommon.StringError(err)
	}

	// Run safety checks
	p, err = t.safetyCheck(ctx, p)
	if err != nil {
		return res, libcommon.StringError(err)
	}

	// Send request to the blockchain and update model status, hash, transaction amount
	p, err = t.initiateTransaction(ctx, p)
	if err != nil {
		return res, libcommon.StringError(err)
	}

	// this Executor will not exist in scope of postProcess
	(*p.executor).Close()

	// Send required information to new thread and return txId to the endpoint. Create a new context since this will run in background
	ctx2 := context.Background()
	go t.postProcess(ctx2, p)

	ids := []string{}
	urls := []string{}
	for i, id := range p.txIds {
		// check if lowercase version of p.executionRequest.Quote.TransactionRequest.Actions[i].CxFunc contains the word "approve"
		if !strings.Contains(strings.ToLower(p.executionRequest.Quote.TransactionRequest.Actions[i].CxFunc), "approve") {
			ids = append(ids, id)
			urls = append(urls, p.chain.Explorer+"/tx/"+id)
		}
	}
	return model.TransactionReceipt{TxIds: ids, TxURLs: urls, TxTimestamp: time.Now().Format(time.RFC1123)}, nil
}

func (t transaction) transactionSetup(ctx context.Context, p transactionProcessingData) (transactionProcessingData, error) {
	_, finish := Span(ctx, "service.transaction.transactionSetup", SpanTag{"platformId": p.platformId})
	defer finish()

	user, err := t.repos.User.GetById(ctx, *p.userId)
	if err != nil {
		return p, libcommon.StringError(err)
	}
	email, err := t.repos.Contact.GetByUserIdAndType(ctx, user.Id, "email")
	if err != nil && errors.Cause(err).Error() != "not found" {
		return p, libcommon.StringError(err)
	}
	user.Email = email.Data
	p.user = &user

	// Pull chain info needed for execution from repository
	chain, err := ChainInfo(ctx, p.executionRequest.Quote.TransactionRequest.ChainId, t.repos.Network, t.repos.Asset)
	if err != nil {
		return p, libcommon.StringError(err)
	}
	p.chain = &chain

	// Create new Tx in repository, populate it with known info
	transactionModel, err := t.repos.Transaction.Create(ctx, model.Transaction{Status: "Created", NetworkId: chain.UUID, DeviceId: *p.deviceId, IPAddress: *p.ip, PlatformId: *p.platformId})
	if err != nil {
		return p, libcommon.StringError(err)
	}
	p.transactionModel = &transactionModel

	updateDB := &model.TransactionUpdates{}
	processingFeeAsset, err := t.populateInitialTxModelData(ctx, *p.executionRequest, updateDB)
	p.processingFeeAsset = &processingFeeAsset
	if err != nil {
		return p, libcommon.StringError(err)
	}
	err = t.repos.Transaction.Update(ctx, transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Send()
		return p, libcommon.StringError(err)
	}

	// Dial the RPC and update model status
	executor := NewExecutor()
	p.executor = &executor
	err = executor.Initialize(chain)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	err = t.updateTransactionStatus(ctx, "RPC Dialed", transactionModel.Id)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	return p, err
}

func (t transaction) safetyCheck(ctx context.Context, p transactionProcessingData) (transactionProcessingData, error) {
	_, finish := Span(ctx, "service.transaction.safetyCheck", SpanTag{"platformId": p.platformId})
	defer finish()

	// Test the Tx and update model status
	estimateUSD, estimateETH, estimateEVM, err := t.testTransaction(*p.executor, p.executionRequest.Quote.TransactionRequest, *p.chain, false, false)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	err = t.updateTransactionStatus(ctx, "Tested and Estimated", p.transactionModel.Id)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	// Verify the Quote and update model status
	_, err = verifyQuote(*p.executionRequest, estimateUSD)
	if err != nil {
		// Update cache if price is too volatile
		if errors.Cause(err).Error() == "verifyQuote: price too volatile" {
			quoteCache := NewQuoteCache(t.redis)
			err = quoteCache.PutCachedTransactionRequest(p.executionRequest.Quote.TransactionRequest, estimateEVM)
			if err != nil {
				return p, libcommon.StringError(err)
			}
		}
		return p, libcommon.StringError(err)
	}

	err = t.updateTransactionStatus(ctx, "Quote Verified", p.transactionModel.Id)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	floatEstimate := common.EstimateToImprecise(p.executionRequest.Quote.Estimate)
	p.floatEstimate = &floatEstimate

	// Get current balance of primary token
	balance, err := (*p.executor).GetBalance()
	if err != nil {
		return p, libcommon.StringError(err)
	}

	// Notify staff if balance is below threshold
	threshold := 1.0
	if balance-estimateETH < threshold {
		msg := fmt.Sprintf("STRING-API: %s balance is at or below threshold of %.2f before executing %.2f transaction at %.2f", p.chain.OwlracleName, threshold, estimateETH, balance)
		go MessageTeam(msg)
	}

	// Exit if transaction will fail due to insufficient balance
	if balance <= estimateETH {
		return p, libcommon.StringError(errors.New("hot wallet ETH balance too low"))
	}

	// Authorize quoted cost on end-user CC and update model status
	p, err = t.authCard(ctx, p)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	// Validate Transaction through Real Time Rules engine
	txModel, err := t.repos.Transaction.GetById(ctx, p.transactionModel.Id)
	if err != nil {
		log.Err(err).Msg("error getting tx model in unit21 Tx Evalute")
		return p, libcommon.StringError(err)
	}

	results, err := t.unit21.Transaction.Evaluate(ctx, txModel)
	if err != nil {
		// If Unit21 Evaluate fails, just log, but otherwise continue with the transaction
		log.Err(err).Msg("Error evaluating transaction in Unit21")
		return p, nil // NOTE: intentionally returning nil here in order to continue the transaction
	}

	if len(results) > 0 {
		err = t.updateTransactionStatus(ctx, "Failed", p.transactionModel.Id)
		if err != nil {
			return p, libcommon.StringError(err)
		}

		err = t.unit21CreateTransaction(ctx, p.transactionModel.Id)
		if err != nil {
			return p, libcommon.StringError(err)
		}
		errorString := "\n"
		for _, rule := range results {
			errorString += fmt.Sprintln("Rule: ", rule.RuleName, " - ", rule.Status)
		}
		return p, libcommon.StringError(errors.New(fmt.Sprintf("risk: Transaction Failed Unit21 Real Time Rules Evaluation with results: %+v", errorString)))
	}

	err = t.updateTransactionStatus(ctx, "Unit21 Authorized", p.transactionModel.Id)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	return p, nil
}

func (t transaction) initiateTransaction(ctx context.Context, p transactionProcessingData) (transactionProcessingData, error) {
	_, finish := Span(ctx, "service.transaction.initiateTransaction", SpanTag{"platformId": p.platformId})
	defer finish()

	request := p.executionRequest.Quote.TransactionRequest
	calls := []ContractCall{}
	for _, action := range request.Actions {

		call := ContractCall{
			CxAddr:     action.CxAddr,
			CxFunc:     action.CxFunc,
			CxReturn:   action.CxReturn,
			CxParams:   action.CxParams,
			TxValue:    action.TxValue,
			TxGasLimit: action.TxGasLimit,
		}
		calls = append(calls, call)
	}

	txIds, value, err := (*p.executor).Initiate(calls)
	p.cumulativeValue = value
	if err != nil {
		return p, libcommon.StringError(err)
	}
	p.txIds = append(p.txIds, txIds...)

	// Create Response Tx leg
	eth := common.WeiToEther(value)
	wei := floatToFixedString(eth, 18)
	usd := floatToFixedString(p.floatEstimate.TotalUSD, int(p.processingFeeAsset.Decimals))
	responseLeg := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetId:      p.processingFeeAsset.Id,
		UserId:       *p.userId,
		InstrumentId: t.ids.StringWalletId,
	}
	responseLeg, err = t.repos.TxLeg.Create(ctx, responseLeg)
	if err != nil {
		return p, libcommon.StringError(err)
	}
	txLeg := model.TransactionUpdates{ResponseTxLegId: &responseLeg.Id}
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, txLeg)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	status := "Transaction Initiated"
	txAmount := p.cumulativeValue.String()

	hashesOtherThanApprove := []string{}
	for i, txId := range p.txIds {
		if !strings.Contains(strings.ToLower(p.executionRequest.Quote.TransactionRequest.Actions[i].CxFunc), "approve") {
			hashesOtherThanApprove = append(hashesOtherThanApprove, txId)
		}
	}

	hashes := strings.Join(hashesOtherThanApprove, ", ")
	updateDB := &model.TransactionUpdates{Status: &status, TransactionHash: &hashes, TransactionAmount: &txAmount}
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, updateDB)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	return p, nil
}

func (t transaction) postProcess(ctx context.Context, p transactionProcessingData) {
	_, finish := Span(ctx, "service.transaction.postProcess", SpanTag{"platformId": p.platformId})
	defer finish()

	// Reinitialize Executor
	executor := NewExecutor()
	p.executor = &executor
	err := executor.Initialize(*p.chain)
	if err != nil {
		log.Err(err).Msg("Failed to initialized executor in postProcess")
		// TODO: Handle error instead of returning it
	}

	// Update TX Status
	updateDB := model.TransactionUpdates{}
	status := "Post Process RPC Dialed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Post Process RPC Dialed'")
		// TODO: Handle error instead of returning it
	}

	// confirm the Tx on the EVM
	trueGas, err := confirmTx(executor, p.txIds)
	p.trueGas = trueGas
	if err != nil {
		log.Err(err).Msg("Failed to confirm transaction")
		// TODO: Handle error instead of returning it
	}

	// Update DB status and NetworkFee
	status = "Tx Confirmed"
	updateDB.Status = &status
	networkFee := strconv.FormatUint(trueGas, 10)
	updateDB.NetworkFee = &networkFee // geth uses uint64 for gas
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Tx Confirmed'")
		// TODO: Handle error instead of returning it
	}

	// Get the Token IDs which were transferred
	tokenIds, err := executor.GetTokenIds(p.txIds)
	if err != nil {
		log.Err(err).Msg("Failed to get token ids")
		// TODO: Handle error instead of returning it
	}
	p.tokenIds = strings.Join(tokenIds, ",")

	// Forward any non fungible tokens received to the user
	// TODO: Use the TX ID/s from this in the receipt
	// TODO: Find a way to charge for the gas used in this transaction
	if len(tokenIds) > 0 {
		forwardTxIds /*forwardTokenIds*/, _, err := executor.ForwardNonFungibleTokens(p.txIds, p.executionRequest.Quote.TransactionRequest.UserAddress)
		if err != nil {
			log.Err(err).Msg("Failed to forward non fungible tokens")
		}
		p.forwardTxIds = append(p.forwardTxIds, forwardTxIds...)
	}

	// Get the Token quantities which were transferred
	tokenQuantities, err := executor.GetTokenQuantities(p.txIds)
	if err != nil {
		log.Err(err).Msg("Failed to get token quantities")
		// TODO: Handle error instead of returning it
	}
	p.tokenQuantities = strings.Join(tokenQuantities, ",")

	if len(tokenQuantities) > 0 {
		forwardTxIds /*forwardTokenAddresses*/, _ /*tokenQuantities*/, _, err := executor.ForwardTokens(p.txIds, p.executionRequest.Quote.TransactionRequest.UserAddress)
		if err != nil {
			log.Err(err).Msg("Failed to forward tokens")
		}
		p.forwardTxIds = append(p.forwardTxIds, forwardTxIds...)
	}

	// Cull any TXIDs from Approve(), the user and Unit21 and the receipt don't need them
	for i, action := range p.executionRequest.Quote.TransactionRequest.Actions {
		if strings.Contains(strings.ToLower(action.CxFunc), "approve") {
			p.txIds = append(p.txIds[:i], p.txIds[i+1:]...)
		}
	}

	// TODO: Get the final gas total here and cache it to the quote cache.  And use it for subsequent quotes.
	forwardGas, err := confirmTx(executor, p.forwardTxIds)
	if err != nil {
		log.Err(err).Msg("Failed to confirm forwarding transactions")
	}
	p.trueGas += forwardGas

	// We can close the executor because we aren't using it after this
	executor.Close()

	// Cache the gas associated with this transaction
	qc := NewQuoteCache(t.redis)
	err = qc.UpdateMaxCachedTrueGas(p.executionRequest.Quote.TransactionRequest, p.trueGas)
	if err != nil {
		log.Err(err).Msg("Failed to update quote true gas cache")
	}

	// compute profit
	// TODO: factor request.processingFeeAsset in the event of crypto-to-usd
	profit, err := t.tenderTransaction(ctx, p)
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
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Profit Tendered'")
		// TODO: Handle error instead of returning it
	}

	// charge the users CC
	err = t.chargeCard(ctx, p)
	if err != nil {
		log.Err(err).Msg("failed to charge card")
		// TODO: Handle error instead of returning it
	}

	// Update status upon success
	status = "Card Charged"
	updateDB.Status = &status
	// TODO: Figure out how much we paid the CC payment processor and deduct it
	// and use it to populate processing_fee and processing_fee_asset in the table
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Card Charged'")
		// TODO: Handle error instead of returning it
	}

	// Transaction complete!  Update status
	status = "Completed"
	updateDB.Status = &status
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, updateDB)
	if err != nil {
		log.Err(err).Msg("Failed to update transaction repo with status 'Completed'")
	}

	// Create Transaction data in Unit21
	err = t.unit21CreateTransaction(ctx, p.transactionModel.Id)
	if err != nil {
		log.Err(err).Msg("Error creating Unit21 transaction")
	}

	// send email receipt
	err = t.sendEmailReceipt(ctx, p)
	if err != nil {
		log.Err(err).Msg("Error sending email receipt to user")
	}
}

func (t transaction) populateInitialTxModelData(ctx context.Context, e model.ExecutionRequest, m *model.TransactionUpdates) (model.Asset, error) {
	txType := "fiat-to-crypto"
	m.Type = &txType
	// TODO populate transactionModel.Tags with key-val pairs for Unit21
	// TODO populate transactionModel.DeviceId with info from fingerprint
	// TODO populate transactionModel.IPAddress with info from fingerprint
	// TODO populate transactionModel.PlatformId with UUID of customer
	// bytes, err := json.Marshal()

	// For now just concat everything
	concatParams := []string{}
	concatFuncs := []string{}
	for _, action := range e.Quote.TransactionRequest.Actions {
		concatParams = append(concatParams, "[")
		concatParams = append(concatParams, action.CxParams...)
		concatParams = append(concatParams, "]")
		concatFuncs = append(concatFuncs, action.CxFunc+action.CxReturn)
	}

	contractParams := pq.StringArray(concatParams)
	m.ContractParams = &contractParams
	contractFunc := strings.Join(concatFuncs, ",")
	m.ContractFunc = &contractFunc

	asset, err := t.repos.Asset.GetByName(ctx, "USD")
	if err != nil {
		return model.Asset{}, libcommon.StringError(err)
	}
	m.ProcessingFeeAsset = &asset.Id // Checkout processing asset
	return asset, nil
}

func (t transaction) testTransaction(executor Executor, request model.TransactionRequest, chain Chain, useBuffer bool, useCache bool) (model.Estimate[float64], float64, CallEstimate, error) {
	res := model.Estimate[float64]{}

	quoteCache := NewQuoteCache(t.redis)
	estimateEVM := CallEstimate{}
	recalculate := true
	var err error
	if useBuffer {
		recalculate, estimateEVM, err = quoteCache.CheckUpdateCachedTransactionRequest(request, 60*5) // TODO: robust buffer time
		if err != nil {
			return res, 0, CallEstimate{}, libcommon.StringError(err)
		}
	}

	if recalculate {
		calls := []ContractCall{}
		for _, action := range request.Actions {
			calls = append(calls, ContractCall{
				CxAddr:     action.CxAddr,
				CxFunc:     action.CxFunc,
				CxReturn:   action.CxReturn,
				CxParams:   action.CxParams,
				TxValue:    action.TxValue,
				TxGasLimit: action.TxGasLimit,
			})
		}
		// Estimate value and gas of Tx request
		estimateEVM, err = executor.Estimate(calls)
		if err != nil {
			return res, 0, CallEstimate{}, libcommon.StringError(err)
		}
		if useCache {
			err = quoteCache.PutCachedTransactionRequest(request, estimateEVM)
			if err != nil {
				return res, 0, CallEstimate{}, libcommon.StringError(err)
			}
		}
	}

	// Factor in approvals to Token Cost
	tokenAddresses := []string{}
	tokenAmounts := []big.Int{}
	for _, action := range request.Actions {
		if strings.ToLower(strings.ReplaceAll(action.CxFunc, " ", "")) == "approve(address,uint256)" {
			tokenAddresses = append(tokenAddresses, action.CxAddr)
			// It should be safe at this point to w3.I without panic
			tokenAmounts = append(tokenAmounts, *w3.I(action.CxParams[1]))
		}
	}

	// Calculate total eth estimate as float64
	gas := new(big.Int)
	gas.SetUint64(estimateEVM.Gas)
	wei := gas.Add(&estimateEVM.Value, gas)
	eth := common.WeiToEther(wei)

	chainId, err := executor.GetByChainId()
	if err != nil {
		return res, eth, CallEstimate{}, libcommon.StringError(err)
	}
	cost := NewCost(t.redis)
	estimationParams := EstimationParams{
		ChainId:    chainId,
		CostETH:    estimateEVM.Value,
		UseBuffer:  useBuffer,
		GasUsedWei: estimateEVM.Gas,
		CostTokens: tokenAmounts,
		TokenAddrs: tokenAddresses,
	}

	// Estimate Cost in USD to execute Tx request
	estimateUSD, err := cost.EstimateTransaction(estimationParams, chain)
	if err != nil {
		return res, eth, CallEstimate{}, libcommon.StringError(err)
	}
	res = estimateUSD
	return res, eth, estimateEVM, nil
}

func verifyQuote(e model.ExecutionRequest, newEstimate model.Estimate[float64]) (bool, error) {
	// Null out values which have changed since payload was signed
	dataToValidate := e
	dataToValidate.Quote.Signature = ""
	bytesToValidate, err := json.Marshal(dataToValidate.Quote)
	if err != nil {
		return false, libcommon.StringError(err)
	}
	valid, err := common.ValidateEVMSignature(e.Quote.Signature, bytesToValidate, true)
	if err != nil {
		return false, libcommon.StringError(err)
	}
	if !valid {
		return false, libcommon.StringError(errors.New("verifyQuote: invalid signature"))
	}
	if newEstimate.Timestamp-e.Quote.Estimate.Timestamp > 20 {
		return false, libcommon.StringError(errors.New("verifyQuote: quote expired"))
	}
	quotedTotal, err := strconv.ParseFloat(e.Quote.Estimate.TotalUSD, 64)
	if err != nil {
		return false, libcommon.StringError(err)
	}
	if newEstimate.TotalUSD > quotedTotal {
		return false, libcommon.StringError(errors.New("verifyQuote: price too volatile"))
	}
	return true, nil
}

func (t transaction) addCardInstrumentIdIfNew(ctx context.Context, p transactionProcessingData) (string, error) {
	_, finish := Span(ctx, "service.transaction.addCardInstrumentIdIfNew", SpanTag{"platformId": p.platformId})
	defer finish()
	// Create a new context since there are sub routines that run in background
	ctx2 := context.Background()

	instrument, err := t.repos.Instrument.GetCardByFingerprint(ctx, p.cardAuthorization.CheckoutFingerprint)
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return "", libcommon.StringError(err)
	} else if err == nil && instrument.UserId != "" {
		go t.unit21.Instrument.Update(ctx2, instrument) // if instrument already exists, update it anyways
		return instrument.Id, nil                       // return if instrument already exists
	}

	// We should gather type from the payment processor
	instrument_type := "debit card"
	if p.cardAuthorization.CardType == "CREDIT" {
		instrument_type = "credit card"
	}
	// Create a new instrument
	instrument = model.Instrument{ // No locationId until fingerprint
		Type:      instrument_type,
		Status:    "created",
		Last4:     p.cardAuthorization.Last4,
		UserId:    *p.userId,
		PublicKey: p.cardAuthorization.CheckoutFingerprint,
		Name:      p.cardAuthorization.CardholderName,
	}

	instrument, err = t.repos.Instrument.Create(ctx, instrument)
	if err != nil {
		return "", libcommon.StringError(err)
	}

	go t.unit21.Instrument.Create(ctx2, instrument)

	return instrument.Id, nil
}

func (t transaction) addWalletInstrumentIdIfNew(ctx context.Context, address string, id string) (string, error) {
	_, finish := Span(ctx, "service.transaction.addWalletInstrumentIdIfNew")
	defer finish()

	// Create a new context since this will run in background
	ctx2 := context.Background()

	instrument, err := t.repos.Instrument.GetWalletByAddr(ctx, address)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return "", libcommon.StringError(err)
	} else if err == nil && instrument.PublicKey == address {
		go t.unit21.Instrument.Update(ctx2, instrument) // if instrument already exists, update it anyways
		return instrument.Id, nil                       // return if instrument already exists
	}

	// Create a new instrument
	instrument = model.Instrument{Type: "crypto wallet", Status: "external", Network: "ethereum", PublicKey: address, UserId: id} // No locationId or userId because this wallet was not registered with the user and is some other recipient
	instrument, err = t.repos.Instrument.Create(ctx, instrument)
	if err != nil {
		return "", libcommon.StringError(err)
	}

	go t.unit21.Instrument.Create(ctx2, instrument)

	return instrument.Id, nil
}

func (t transaction) authCard(ctx context.Context, p transactionProcessingData) (transactionProcessingData, error) {
	_, finish := Span(ctx, "service.transaction.authCard", SpanTag{"platformId": p.platformId})
	defer finish()

	// auth their card
	p, err := AuthorizeCharge(p)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	// Add Checkout Instrument ID to our DB if it's not there already and associate it with the user
	instrumentId, err := t.addCardInstrumentIdIfNew(ctx, p)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	// Create Origin Tx leg
	usdWei := floatToFixedString(p.floatEstimate.TotalUSD, int(p.processingFeeAsset.Decimals))
	origin := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetId:      p.processingFeeAsset.Id,
		UserId:       *p.userId,
		InstrumentId: instrumentId,
	}
	origin, err = t.repos.TxLeg.Create(ctx, origin)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	txLegUpdates := model.TransactionUpdates{OriginTxLegId: &origin.Id}
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, txLegUpdates)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	err = t.updateTransactionStatus(ctx, "Card "+p.cardAuthorization.Status, p.transactionModel.Id)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	recipientWalletId, err := t.addWalletInstrumentIdIfNew(ctx, p.executionRequest.Quote.TransactionRequest.UserAddress, *p.userId)
	p.recipientWalletId = &recipientWalletId
	if err != nil {
		return p, libcommon.StringError(err)
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

	destinationLeg, err = t.repos.TxLeg.Create(ctx, destinationLeg)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	txLegUpdates = model.TransactionUpdates{DestinationTxLegId: &destinationLeg.Id}

	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, txLegUpdates)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	if !p.cardAuthorization.Approved {
		go t.unit21CreateTransaction(ctx, p.transactionModel.Id)

		return p, libcommon.StringError(errors.New("payment: Authorization Declined by Checkout"))
	}

	return p, nil
}

func confirmTx(executor Executor, txIds []string) (uint64, error) {
	trueGas, err := executor.TxWait(txIds)
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	return trueGas, nil
}

// TODO: rewrite this transaction to reference the asset(s) received by the user, not what we paid
func (t transaction) tenderTransaction(ctx context.Context, p transactionProcessingData) (float64, error) {
	_, finish := Span(ctx, "service.transaction.tenderTransaction", SpanTag{"platformId": p.platformId})
	defer finish()

	cost := NewCost(t.redis)
	trueWei := big.NewInt(0).Add(p.cumulativeValue, big.NewInt(int64(p.trueGas)))
	trueEth := common.WeiToEther(trueWei)
	trueUSD, err := cost.LookupUSD(trueEth, p.chain.CoingeckoName, p.chain.CoincapName)
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	profit := p.floatEstimate.TotalUSD - trueUSD

	// Create Receive Tx leg
	asset, err := t.repos.Asset.GetById(ctx, p.chain.GasTokenId)
	if err != nil {
		return profit, libcommon.StringError(err)
	}
	wei := floatToFixedString(trueEth, int(asset.Decimals))
	usd := floatToFixedString(p.floatEstimate.TotalUSD, 6)

	txModel, err := t.repos.Transaction.GetById(ctx, p.transactionModel.Id)
	if err != nil {
		return profit, libcommon.StringError(err)
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
	err = t.repos.TxLeg.Update(ctx, txModel.DestinationTxLegId, destinationLeg)
	if err != nil {
		return profit, libcommon.StringError(err)
	}

	return profit, nil
}

func (t transaction) chargeCard(ctx context.Context, p transactionProcessingData) error {
	_, finish := Span(ctx, "service.transaction.chargeCard", SpanTag{"platformId": p.platformId})
	defer finish()

	p, err := CaptureCharge(p)
	if err != nil {
		return libcommon.StringError(err)
	}

	// Create Receipt Tx leg
	usdWei := floatToFixedString(p.floatEstimate.TotalUSD, int(p.processingFeeAsset.Decimals))
	receiptLeg := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetId:      p.processingFeeAsset.Id,
		UserId:       t.ids.StringUserId,
		InstrumentId: t.ids.StringBankId,
	}
	receiptLeg, err = t.repos.TxLeg.Create(ctx, receiptLeg)
	if err != nil {
		return libcommon.StringError(err)
	}
	txLeg := model.TransactionUpdates{ReceiptTxLegId: &receiptLeg.Id, PaymentCode: &p.PaymentId}
	err = t.repos.Transaction.Update(ctx, p.transactionModel.Id, txLeg)
	if err != nil {
		return libcommon.StringError(err)
	}

	return nil
}

func (t transaction) sendEmailReceipt(ctx context.Context, p transactionProcessingData) error {
	_, finish := Span(ctx, "service.transaction.sendEmailReceipt", SpanTag{"platformId": p.platformId})
	defer finish()

	user, err := t.repos.User.GetById(ctx, *p.userId)
	if err != nil {
		log.Err(err).Msg("Error getting user from repo")
		return libcommon.StringError(err)
	}

	contact, err := t.repos.Contact.GetByUserId(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Error getting user contact from repo")
		return libcommon.StringError(err)
	}

	name := user.FirstName // + " " + user.MiddleName + " " + user.LastName
	if name == "" {
		name = "User"
	}

	platform, err := t.repos.Platform.GetById(ctx, *p.platformId)
	if err != nil {
		return libcommon.StringError(err)
	}

	transactionRequest := p.executionRequest.Quote.TransactionRequest
	estimate := p.floatEstimate

	explorers := []string{}
	for _, id := range p.txIds {
		explorers = append(explorers, p.chain.Explorer+"/tx"+id)
	}

	receiptParams := emailer.ReceiptGenerationParams{
		ReceiptType:         "NFT Purchase", // TODO: retrieve dynamically
		CustomerName:        name,
		StringPaymentId:     p.transactionModel.Id,
		PaymentDescriptor:   p.executionRequest.Quote.TransactionRequest.AssetName,
		TransactionDate:     time.Now().Format(time.RFC1123),
		TransactionId:       p.txIds[0],   // For now assume there were 2 and the approval one was removed
		TransactionExplorer: explorers[0], // For now assume there were 2 and the approval one was removed
		DestinationAddress:  transactionRequest.UserAddress,
		DestinationExplorer: p.chain.Explorer + "/address/" + transactionRequest.UserAddress,
		PaymentMethod:       p.cardAuthorization.Issuer + " " + p.cardAuthorization.Last4,
		Platform:            platform.Name,
		ItemOrdered:         p.executionRequest.Quote.TransactionRequest.AssetName,
		TokenId:             p.tokenIds,
		Subtotal:            common.FloatToUSDString(estimate.BaseUSD + estimate.TokenUSD),
		NetworkFee:          common.FloatToUSDString(estimate.GasUSD),
		ProcessingFee:       common.FloatToUSDString(estimate.ServiceUSD),
		Total:               common.FloatToUSDString(estimate.TotalUSD),
	}

	emailer := emailer.New()

	err = emailer.SendReceipt(ctx, contact.Data, receiptParams)
	if err != nil {
		log.Err(err).Msg("Error sending email receipt to user")
		return libcommon.StringError(err)
	}
	return nil
}

func floatToFixedString(value float64, decimals int) string {
	return strconv.FormatUint(uint64(value*(math.Pow10(decimals))), 10)
}

func (t transaction) unit21CreateTransaction(ctx context.Context, transactionId string) (err error) {
	_, finish := Span(ctx, "service.transaction.unit21CreateTransaction", SpanTag{"transactionId": transactionId})
	defer finish()

	txModel, err := t.repos.Transaction.GetById(ctx, transactionId)
	if err != nil {
		log.Err(err).Msg("Error getting tx model in Unit21 in Tx Postprocess")
		return libcommon.StringError(err)
	}

	_, err = t.unit21.Transaction.Create(ctx, txModel)
	if err != nil {
		log.Err(err).Msg("Error updating unit21 in Tx Postprocess")
		return libcommon.StringError(err)
	}

	return nil
}

func (t transaction) updateTransactionStatus(ctx context.Context, status string, transactionId string) (err error) {
	_, finish := Span(ctx, "service.transaction.updateTransactionStatus", SpanTag{"transactionId": transactionId})
	defer finish()

	updateDB := &model.TransactionUpdates{Status: &status}
	err = t.repos.Transaction.Update(ctx, transactionId, updateDB)
	if err != nil {
		return libcommon.StringError(err)
	}

	return nil
}

func (t *transaction) getStringInstrumentsAndUserId() {
	t.ids = GetStringIdsFromEnv()
}

func (t transaction) isContractAllowed(ctx context.Context, platformId string, networkId string, request model.TransactionRequest) (isAllowed bool, err error) {
	_, finish := Span(ctx, "service.transaction.isContractAllowed", SpanTag{"platformId": platformId})
	defer finish()

	for _, action := range request.Actions {
		cxAddr := action.CxAddr
		contract, err := t.repos.Contract.GetForValidation(ctx, cxAddr, networkId, platformId)
		if err != nil && err == serror.NOT_FOUND {
			return false, libcommon.StringError(serror.CONTRACT_NOT_ALLOWED)
		} else if err != nil {
			return false, libcommon.StringError(err)
		}

		if len(contract.Functions) == 0 {
			continue
		}

		for _, function := range contract.Functions {
			if function == action.CxFunc {
				continue
			}
		}
	}

	return true, nil
}
