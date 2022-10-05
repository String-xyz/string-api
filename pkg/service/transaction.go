package service

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Quote(d model.TransactionRequest) (model.ExecutionRequest, error)
	Execute(e model.ExecutionRequest) (model.Transaction, error)
	New(repo repository.Transaction) Transaction
}

type transaction struct {
	repository repository.Transaction
}

func (t transaction) New(repo repository.Transaction) Transaction {
	return &transaction{repository: repo}
}

func NewTransaction(repo repository.Transaction) Transaction {
	return &transaction{repository: repo}
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

func (t transaction) Execute(e model.ExecutionRequest) (model.Transaction, error) {
	res := model.Transaction{}
	// TODO: Create entry of E in TX DB

	chain, err := model.ChainInfo(uint64(e.ChainID))
	if err != nil {
		return res, err
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return res, err
	}

	estimateUSD, err := testTransaction(executor, e.TransactionRequest, false)
	if err != nil {
		return res, err
	}
	// model.status = tested, update db

	_, err = verifyQuote(e, estimateUSD)
	if err != nil {
		return res, err
	}
	// model.status = quoteVerified, update db

	//Authorize quoted cost on end-user CC
	authorizationID, err := authcard(e.UserAddress, e.CardToken, e.TotalUSD)
	if err != nil {
		return res, err
	}
	// model.status = ccAuthorized, update db

	txID, value, err := initiateTransaction(executor, e)
	if err != nil {
		return res, err
	}
	// model.status = txInitiated, update db

	// this Executor will not exist in scope of postProcess
	executor.Close()

	post := postProcessRequest{
		TxID:            txID,
		ChainID:         chain.ChainID,
		AuthorizationID: authorizationID,
		UserAddress:     e.UserAddress,
		CumulativeValue: value,
		QuotedTotal:     e.TotalUSD,
	}
	go postProcess(post)

	return model.Transaction{TxID: txID}, nil
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

	chainID, err := executor.ChainID()
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

func authcard(userWallet string, cardToken string, usd float64) (string, error) {
	// auth their card
	auth, err := Authorize(usd, userWallet, cardToken)
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
	_, err := Capture(usd, userWallet, authorizationID)
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
}

func postProcess(request postProcessRequest) error {
	chain, err := model.ChainInfo(request.ChainID)
	if err != nil {
		return err
	}
	executor := NewExecutor()
	err = executor.Initialize(chain.RPC)
	if err != nil {
		return err
	}
	// confirm the TX on the EVM, update db status
	trueGas, err := confirmTX(executor, request.TxID)
	if err != nil {
		return err
	}

	// compute profit and log to db
	profit, err := tenderTransaction(request.CumulativeValue, trueGas, request.QuotedTotal, chain)
	if err != nil {
		return err
	}
	fmt.Printf("PROFIT=%+v", profit)
	// log profit to db

	// charge the users CC
	err = chargeCard(request.UserAddress, request.AuthorizationID, request.QuotedTotal)
	if err != nil {
		return err
	}

	// update tx status to complete in the db
	executor.Close()
	return nil
}
