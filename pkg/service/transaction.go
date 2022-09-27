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

	estimateUSD, err := testTransaction(d, true)
	if err != nil {
		return res, err
	}
	res.Quote = estimateUSD

	// Sign entire payload
	signature, err := common.Sign(res)
	if err != nil {
		return res, err
	}
	res.Signature = signature

	return res, nil
}

func (t transaction) Execute(e model.ExecutionRequest) (model.Transaction, error) {
	res := model.Transaction{}
	// TODO: Create entry of E in TX DB

	estimateUSD, err := testTransaction(e.TransactionRequest, false)
	if err != nil {
		return res, err
	}
	// model.status = tested, update db
	fmt.Printf("\nestimateUSD=%+v", estimateUSD)

	_, err = verifyQuote(e, estimateUSD)
	if err != nil {
		return res, err
	}
	// model.status = quoteVerified, update db

	//Authorize quoted cost on end-user CC

	// TEST TX
	txID, err := initiateTransaction(e)
	if err != nil {
		fmt.Printf("ERR=%+v", err)
		return res, err
	}
	fmt.Printf("TXID=%+v", txID)

	defer postProcess(txID, e)

	// create execution response struct and return that
	return model.Transaction{TxID: txID}, nil
}

func testTransaction(t model.TransactionRequest, useBuffer bool) (model.Quote, error) {
	res := model.Quote{}
	executor := NewExecutor() // maybe scope this outside and pass in a reference

	// Verify Chain is supported and get RPC for chain
	chain, err := model.ChainInfo(uint64(t.ChainID))
	if err != nil {
		return res, err
	}

	executor.Initialize(chain.RPC) // we want to avoid calling this redundantly in /transact
	call := ContractCall{
		RPC:        chain.RPC,
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
	executor.Close()

	cost := NewCost(repository.NewCost(nil))
	estimationParams := EstimationParams{
		ChainID:     chain.ChainID,
		CostETH:     estimateEVM.Value,
		UseBuffer:   useBuffer,
		GasUsedGwei: estimateEVM.Gas,
		CostToken:   *big.NewInt(0),
		TokenName:   "",
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
	dataToValidate := e
	dataToValidate.Signature = ""
	dataToValidate.CardToken = ""
	valid, err := common.ValidateSignature(e.Signature, dataToValidate)
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

func authcard(userWallet string, cardToken string, usd uint64) error {
	// auth their card
	return nil
}

func initiateTransaction(e model.ExecutionRequest) (string, error) {
	executor := NewExecutor() // maybe scope this outside and pass in a reference
	// Verify Chain is supported and get RPC for chain
	chain, err := model.ChainInfo(uint64(e.ChainID))
	if err != nil {
		return "", err
	}
	executor.Initialize(chain.RPC)
	call := ContractCall{
		RPC:        chain.RPC,
		CxAddr:     e.CxAddr,
		CxFunc:     e.CxFunc,
		CxReturn:   e.CxReturn,
		CxParams:   e.CxParams,
		TxValue:    e.TxValue,
		TxGasLimit: e.TxGasLimit,
	}
	txID, err := executor.Initiate(call)
	if err != nil {
		return "", err
	}
	return txID, nil
}

func postProcess(txID string, e model.ExecutionRequest /*a ChargeResponse, m TransactionModel*/) {

}
