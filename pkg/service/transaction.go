package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Quote(d model.TransactionRequest) (model.ExecutionRequest, error)
	Execute(e model.ExecutionRequest) (model.TransactionReceipt, error)
	New(repo repository.Transaction) Transaction
}

type transaction struct {
	repo repository.Transaction
}

func (t transaction) New(repo repository.Transaction) Transaction {
	return &transaction{repo: repo}
}

func NewTransaction(repo repository.Transaction) Transaction {
	return &transaction{repo: repo}
}

func (t transaction) Quote(d model.TransactionRequest) (model.ExecutionRequest, error) {
	// TODO: use prefab service to parse d and fill out known params
	res := model.ExecutionRequest{TransactionRequest: d}
	chain, err := model.ChainInfo(uint64(d.ChainID))
	if err != nil {
		return res, err
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, err
	}

	estimateUSD, err := testTransaction(executor, d, true)
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
	// TODO: Create entry of E in TX DB
	db, err := t.repo.Create(model.Transaction{Status: "Created"})
	db.ContractABI = e.CxFunc + e.CxReturn
	db.Parameters, err = json.Marshal(e.CxParams)
	if err != nil {
		return res, err
	}
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}

	chain, err := model.ChainInfo(uint64(e.ChainID))
	if err != nil {
		return res, err
	}
	db.NetworkID = chain.OwlracleName
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, err
	}
	db.Status = "RPC Dialed"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}

	estimateUSD, err := testTransaction(executor, e.TransactionRequest, false)
	if err != nil {
		return res, err
	}
	db.Status = "Tested and Estimated"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}

	_, err = verifyQuote(e, estimateUSD)
	if err != nil {
		return res, err
	}
	db.Status = "Quote Verified"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}

	//Authorize quoted cost on end-user CC
	authorizationID, err := authCard(e.UserAddress, e.CardToken, e.TotalUSD)
	if err != nil {
		return res, err
	}
	db.Status = "Card Authorized"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}

	txID, value, err := initiateTransaction(executor, e)
	if err != nil {
		return res, err
	}
	db.Status = "Transaction Initiated"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return res, err
	}

	// this Executor will not exist in scope of postProcess
	executor.Close()

	post := postProcessRequest{
		TxID:            txID,
		ChainID:         chain.ChainID,
		AuthorizationID: authorizationID,
		UserAddress:     e.UserAddress,
		CumulativeValue: value,
		QuotedTotal:     e.TotalUSD,
		db:              db,
	}
	go t.postProcess(post)

	return model.TransactionReceipt{TxID: txID}, nil
}

func testTransaction(executor Executor, t model.TransactionRequest, useBuffer bool) (model.Quote, error) {
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
	estimateUSD, err := cost.EstimateTransaction(estimationParams)
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

func tenderTransaction(cumulativeValue *big.Int, cumulativeGas uint64, quotedTotal float64, chain model.Chain) (float64, error) {
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
	TxID            string
	ChainID         uint64
	AuthorizationID string
	UserAddress     string
	CumulativeGas   uint64
	CumulativeValue *big.Int
	QuotedTotal     float64
	db              model.Transaction
}

func (t transaction) postProcess(request postProcessRequest) error {
	chain, err := model.ChainInfo(request.ChainID)
	if err != nil {
		return err
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return err
	}
	db := request.db
	db.Status = "Post Process RPC Dialed"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return err
	}
	// confirm the TX on the EVM, update db status
	trueGas, err := confirmTX(executor, request.TxID)
	if err != nil {
		return err
	}
	db.Status = "TX Confirmed"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return err
	}

	// compute profit and log to db
	profit, err := tenderTransaction(request.CumulativeValue, trueGas, request.QuotedTotal, chain)
	if err != nil {
		return err
	}
	fmt.Printf("PROFIT=%+v", profit)
	db.Status = "Profit Tendered"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return err
	}

	// charge the users CC
	err = chargeCard(request.UserAddress, request.AuthorizationID, request.QuotedTotal)
	if err != nil {
		return err
	}
	db.Status = "Card Charged"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return err
	}

	db.Status = "Completed"
	err = t.repo.Update(db.ID, db)
	if err != nil {
		return err
	}
	executor.Close()
	return nil
}
