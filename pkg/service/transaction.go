package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/internal/unit21"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/lib/pq"
)

type Transaction interface {
	Quote(d model.TransactionRequest) (model.ExecutionRequest, error)
	Execute(e model.ExecutionRequest, userId string, deviceId string) (model.TransactionReceipt, error)
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
	repos repository.Repositories
	redis store.RedisStore
	ids   InternalIds
}

func NewTransaction(repos repository.Repositories, redis store.RedisStore) Transaction {
	return &transaction{repos: repos, redis: redis}
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

	estimateUSD, _, err := t.testTransaction(executor, d, chain, true)
	if err != nil {
		return res, common.StringError(err)
	}
	res.Quote = estimateUSD
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

func (t transaction) Execute(e model.ExecutionRequest, userId string, deviceId string) (model.TransactionReceipt, error) {
	res := model.TransactionReceipt{}
	t.getStringInstrumentsAndUserId()

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
	db, err := t.repos.Transaction.Create(model.Transaction{Status: "Created", NetworkID: chain.UUID, DeviceID: deviceId, PlatformID: t.ids.StringPlatformId})
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
	estimateUSD, estimateETH, err := t.testTransaction(executor, e.TransactionRequest, chain, false)
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

	status = "Card " + cardAuthorization.Status
	updateDB.Status = &status
	err = t.repos.Transaction.Update(db.ID, updateDB)
	if err != nil {
		return res, common.StringError(err)
	}

	recipientWalletId, err := t.addWalletInstrumentIdIfNew(e.UserAddress, user.ID)
	if err != nil {
		return res, common.StringError(err)
	}

	// TODO: Determine the output of the transaction (destination leg) with Tracers
	destinationLeg := model.TxLeg{
		Timestamp:    time.Now(),        // Required by the db. Should be updated when the tx occurs
		Amount:       "0",               // Required by Unit21. The amount of the asset received by the user
		Value:        "0",               // Default to '0'. The value of the asset received by the user
		AssetID:      chain.GasTokenID,  // Required by the db. the asset received by the user
		UserID:       user.ID,           // the user who received the asset
		InstrumentID: recipientWalletId, // Required by the db. the instrument which received the asset (wallet usually)
	}

	destinationLeg, err = t.repos.TxLeg.Create(destinationLeg)
	if err != nil {
		return res, common.StringError(err)
	}

	txLeg := model.TransactionUpdates{DestinationTxLegID: &destinationLeg.ID}

	err = t.repos.Transaction.Update(db.ID, txLeg)
	if err != nil {
		return res, common.StringError(err)
	}

	if !cardAuthorization.Approved {
		err := t.unit21CreateTransaction(db.ID)
		if err != nil {
			return res, common.StringError(err)
		}

		return res, common.StringError(common.StringError(errors.New("payment: Authorization Declined by Checkout")))
	}

	// Validate Transaction through Real Time Rules engine
	u21auth, err := t.unit21Evaluate(db.ID)
	if err != nil {
		return res, common.StringError(err)
	}

	if !u21auth {
		status = "Failed"
		updateDB.Status = &status
		err = t.repos.Transaction.Update(db.ID, updateDB)
		if err != nil {
			return res, common.StringError(err)
		}

		err = t.unit21CreateTransaction(db.ID)
		if err != nil {
			return res, common.StringError(err)
		}

		return res, common.StringError(errors.New("risk: Transaction Failed Unit21 Real Time Rules Evaluation"))
	}
	status = "Unit21 Authorized"
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
		Authorization:      cardAuthorization,
		UserAddress:        e.UserAddress,
		CumulativeValue:    value,
		Quote:              e.Quote,
		TxDBID:             db.ID,
		processingFeeAsset: processingFeeAsset,
		preBalance:         preBalance,
		userId:             userId,
		recipientWalletId:  recipientWalletId,
	}
	go t.postProcess(post)

	return model.TransactionReceipt{TxID: txID, TxURL: chain.Explorer + "/tx/" + txID}, nil
}

func (t *transaction) getStringInstrumentsAndUserId() {
	t.ids = GetStringIdsFromEnv()
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

	chainID, err := executor.GetChainID()
	if err != nil {
		return res, eth, common.StringError(err)
	}
	cost := NewCost(t.redis)
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

func (t transaction) addWalletInstrumentIdIfNew(address string, id string) (string, error) {
	instrument, err := t.repos.Instrument.GetWallet(address)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return "", common.StringError(err)
	} else if err == nil && instrument.PublicKey == address {
		return instrument.ID, nil
	}

	// Create a new instrument
	instrument = model.Instrument{Type: "crypto-wallet", Status: "external", Network: "ethereum", PublicKey: address, UserID: id} // No locationID or userID because this wallet was not registered with the user and is some other recipient
	instrument, err = t.repos.Instrument.Create(instrument)
	if err != nil {
		return "", common.StringError(err)
	}
	return instrument.ID, nil
}

func (t transaction) authCard(userWallet string, cardToken string, usd float64, chargeAsset model.Asset, dbID string, userId string) (AuthorizedCharge, error) {
	// auth their card
	auth, err := AuthorizeCharge(usd, userWallet, cardToken)
	if err != nil {
		return auth, common.StringError(err)
	}

	// Add Checkout Instrument ID to our DB if it's not there already and associate it with the user
	instrumentId, err := t.addCardInstrumentIdIfNew(auth.CheckoutFingerprint, userId, auth.Last4)
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

	go t.unit21CreateInstrument(origin.InstrumentID)

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

	// Create Response Tx leg
	eth := common.WeiToEther(value)
	wei := floatToFixedString(eth, 18)
	usd := floatToFixedString(e.TotalUSD, int(chargeAsset.Decimals))
	responseLeg := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       wei,
		Value:        usd,
		AssetID:      chargeAsset.ID,
		UserID:       userId,
		InstrumentID: t.ids.StringWalletId,
	}
	responseLeg, err = t.repos.TxLeg.Create(responseLeg)
	if err != nil {
		return txID, value, common.StringError(err)
	}
	txLeg := model.TransactionUpdates{ResponseTxLegID: &responseLeg.ID}
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
	receiptLeg := model.TxLeg{
		Timestamp:    time.Now(),
		Amount:       usdWei,
		Value:        usdWei,
		AssetID:      chargeAsset.ID,
		UserID:       t.ids.StringUserId,
		InstrumentID: t.ids.StringBankId,
	}
	receiptLeg, err = t.repos.TxLeg.Create(receiptLeg)
	if err != nil {
		return common.StringError(err)
	}
	txLeg := model.TransactionUpdates{ReceiptTxLegID: &receiptLeg.ID}
	err = t.repos.Transaction.Update(txUUID, txLeg)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

// TODO: rewrite this transaction to reference the asset(s) received by the user, not what we paid
func (t transaction) tenderTransaction(cumulativeValue *big.Int, cumulativeGas uint64, quotedTotal float64, chain Chain, txUUID string, recipientId string, userWalletId string) (float64, error) {
	cost := NewCost(t.redis)
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

	txModel, err := t.repos.Transaction.GetById(txUUID)
	if err != nil {
		return profit, common.StringError(err)
	}

	destinationLeg := model.TxLeg{
		Timestamp:    time.Now(),   // updated based on *when the transaction occured* not time.Now()
		Amount:       wei,          // Should be the amount of the asset received by the user
		Value:        usd,          // The value of the asset received by the user
		AssetID:      asset.ID,     // the asset received by the user
		UserID:       recipientId,  // the user who received the asset
		InstrumentID: userWalletId, // the instrument which received the asset (wallet usually)
	}

	// We now update the destination leg instead of creating it
	err = t.repos.TxLeg.Update(txModel.DestinationTxLegID, destinationLeg)
	if err != nil {
		return profit, common.StringError(err)
	}

	return profit, nil
}

type postProcessRequest struct {
	TxID               string
	Chain              Chain
	Authorization      AuthorizedCharge
	UserAddress        string
	CumulativeGas      uint64
	CumulativeValue    *big.Int
	Quote              model.Quote
	TxDBID             string
	processingFeeAsset model.Asset
	preBalance         float64
	userId             string
	recipientWalletId  string
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
	profit, err := t.tenderTransaction(request.CumulativeValue, trueGas, request.Quote.TotalUSD, request.Chain, request.TxDBID, request.userId, request.recipientWalletId)
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
	err = t.chargeCard(request.UserAddress, request.Authorization.AuthID, request.Quote.TotalUSD, request.processingFeeAsset, request.TxDBID, request.userId)
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

	err = t.unit21CreateTransaction(request.TxDBID)
	if err != nil {
		log.Printf("Error creating Unit21 transaction: %s", err)
	}

	// send email receipt
	err = t.sendEmailReceipt(request)
	if err != nil {
		log.Printf("Error sending email receipt to user: %s", err)
	}
}

func (t transaction) sendEmailReceipt(request postProcessRequest) error {
	user, err := t.repos.User.GetById(request.userId)
	if err != nil {
		log.Printf("Error getting user from repo: %s", err)
		return err
	}
	contact, err := t.repos.Contact.GetByUserId(request.userId)
	if err != nil {
		log.Printf("Error getting user contact from repo: %s", err)
		return err
	}
	name := user.FirstName // + " " + user.MiddleName + " " + user.LastName
	if name == "" {
		name = "User"
	}
	receiptParams := common.ReceiptGenerationParams{
		ReceiptType:       "NFT Purchase", // TODO: retrieve dynamically
		CustomerName:      name,
		StringPaymentId:   request.TxDBID,
		PaymentDescriptor: "String Digital Asset", // TODO: retrieve dynamically
		TransactionDate:   time.Now().Format(time.RFC1123),
	}
	receiptBody := [][2]string{
		{"Transaction ID", "<a href='" + request.Chain.Explorer + "/tx/" + request.TxID + "'>" + request.TxID + "</a>"},
		{"Destination Wallet", "<a href='" + request.Chain.Explorer + "/address/" + request.UserAddress + "'>" + request.UserAddress + "</a>"},
		{"Payment Descriptor", receiptParams.PaymentDescriptor},
		{"Payment Method", request.Authorization.Issuer + " " + request.Authorization.Last4},
		{"Platform", "String Demo"},            // TODO: retrieve dynamically
		{"Item Ordered", "String Fighter NFT"}, // TODO: retrieve dynamically
		{"Token ID", "1234"},                   // TODO: retrieve dynamically, maybe after building token transfer detection
		{"Subtotal", common.FloatToUSDString(request.Quote.BaseUSD + request.Quote.TokenUSD)},
		{"Network Fee:", common.FloatToUSDString(request.Quote.GasUSD)},
		{"Processing Fee", common.FloatToUSDString(request.Quote.ServiceUSD)},
		{"Total Charge", common.FloatToUSDString(request.Quote.TotalUSD)},
	}
	err = common.EmailReceipt(contact.Data, receiptParams, receiptBody)
	if err != nil {
		log.Printf("Error sending email receipt to user: %s", err)
		return err
	}
	return nil
}

func floatToFixedString(value float64, decimals int) string {
	return strconv.FormatUint(uint64(value*(math.Pow10(decimals-1))), 10)
}

func (t transaction) unit21CreateInstrument(instrumentId string) (err error) {
	// Send Instrument Data to Unit21
	instrument, err := t.repos.Instrument.GetById(instrumentId)
	if err != nil {
		fmt.Printf("Error creating new instrument in Unit21 -- can't get instrument model")
		return
	}

	u21Repo := unit21.InstrumentRepo{
		User:     t.repos.User,
		Device:   t.repos.Device,
		Location: t.repos.Location, // empty until fingerprint integration
	}

	u21Tx := unit21.NewInstrument(u21Repo)
	_, err = u21Tx.Create(instrument)
	if err != nil {
		fmt.Printf("Error creating new instrument in Unit21")
		return
	}

	return
}

func (t transaction) unit21CreateTransaction(transactionId string) (err error) {
	txModel, err := t.repos.Transaction.GetById(transactionId)
	if err != nil {
		log.Printf("Error getting tx model in Unit21 in Tx Postprocess: %s", err)
		return
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
		return
	}

	return
}

func (t transaction) unit21Evaluate(transactionId string) (evaluation bool, err error) {
	//Check transaction in Unit21
	txModel, err := t.repos.Transaction.GetById(transactionId)
	if err != nil {
		log.Printf("Error getting tx model in Unit21 in Tx Evaluate: %s", err)
		return evaluation, common.StringError(err)
	}

	u21Repo := unit21.TransactionRepo{
		TxLeg: t.repos.TxLeg,
		User:  t.repos.User,
		Asset: t.repos.Asset,
	}

	u21Tx := unit21.NewTransaction(u21Repo)
	evaluation, err = u21Tx.Evaluate(txModel)
	if err != nil {
		log.Printf("Error evaluating transaction in Unit21: %s", err)
		return evaluation, common.StringError(err)
	}

	return

}
