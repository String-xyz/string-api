package service

import (
	"fmt"
	"math/big"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Quote(d model.TransactionData) (model.TransactionRequest, error)
	Execute() (model.Transaction, error)
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

func (t transaction) Quote(d model.TransactionData) (model.TransactionRequest, error) {
	executor := NewExecutor()
	chain, err := model.ChainInfo(uint64(d.ChainID))
	if err != nil {
		return model.TransactionRequest{}, err
	}
	executor.Initialize(chain.RPC)
	call := ContractCall{
		RPC:        chain.RPC,
		CxAddr:     d.ContractAddress,
		CxFunc:     d.ContractFunction,
		CxReturn:   d.ContractReturn,
		CxParams:   d.ContractParameters,
		TxValue:    d.TxValue,
		TxGasLimit: d.GasLimit,
	}
	estimate, err := executor.Estimate(call)
	fmt.Printf("ESTIMATE=%+v\n\n", estimate)
	fmt.Printf("ERR=%+v", err)

	executor.Close()

	cost := NewCost(repository.NewCost(nil))
	estimationParams := EstimationParams{
		ChainID:     chain.ChainID,
		CostETH:     estimate.Value,
		UseBuffer:   true,
		GasUsedGwei: estimate.Gas,
		CostToken:   *big.NewInt(0),
		TokenName:   "",
	}
	res, err := cost.EstimateTransaction(estimationParams)
	if err != nil {
		return model.TransactionRequest{}, err
	}
	fmt.Println("ETH=", res)

	//TEST
	signed, err := common.Sign("this is a test")
	if err != nil {
		return model.TransactionRequest{}, err
	}
	fmt.Printf("\nSIGNED=%+v", signed)
	valid, err := common.ValidateSignature(signed, "this is a test2")
	if err != nil {
		fmt.Printf("\nERROR=%+v", err)
		return model.TransactionRequest{}, err
	}
	fmt.Printf("\nVALID=%+v", valid)

	return model.TransactionRequest{}, nil
}

func (t transaction) Execute() (model.Transaction, error) {
	// create execution response struct and return that
	return model.Transaction{}, nil
}
