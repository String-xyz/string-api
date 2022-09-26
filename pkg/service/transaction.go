package service

import (
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

	executor := NewExecutor()
	// Verify Chain is supported and get RPC for chain
	chain, err := model.ChainInfo(uint64(d.ChainID))
	if err != nil {
		return res, err
	}

	executor.Initialize(chain.RPC)
	call := ContractCall{
		RPC:        chain.RPC,
		CxAddr:     d.CxAddr,
		CxFunc:     d.CxFunc,
		CxReturn:   d.CxReturn,
		CxParams:   d.CxParams,
		TxValue:    d.TxValue,
		TxGasLimit: d.TxGasLimit,
	}
	// Estimate value and gas of TX request
	callEstimate, err := executor.Estimate(call)
	if err != nil {
		return res, err
	}

	executor.Close()

	cost := NewCost(repository.NewCost(nil))
	estimationParams := EstimationParams{
		ChainID:     chain.ChainID,
		CostETH:     callEstimate.Value,
		UseBuffer:   true,
		GasUsedGwei: callEstimate.Gas,
		CostToken:   *big.NewInt(0),
		TokenName:   "",
	}
	// Estimate Cost in USD to execute TX request
	costEstimate, err := cost.EstimateTransaction(estimationParams)
	if err != nil {
		return res, err
	}
	res.Quote = costEstimate

	// Sign entire payload
	signature, err := common.Sign(res)
	if err != nil {
		return res, err
	}
	res.Signature = signature

	return res, nil
}

func (t transaction) Execute(e model.ExecutionRequest) (model.Transaction, error) {
	// TODO: Create entry of E in TX DB

	// create execution response struct and return that
	return model.Transaction{}, nil
}
