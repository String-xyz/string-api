package service

import (
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
	// create quote struct and return that
	//return model.SignedQuote{}, nil
	// cost := NewCost(repository.NewCost(nil))
	// res, err := cost.CoingeckoUSD("ethereum", 1)
	// if err != nil {
	// 	return model.TransactionRequest{}, nil
	// }
	// fmt.Println("ETH=", res)
	return model.TransactionRequest{}, nil
}

func (t transaction) Execute() (model.Transaction, error) {
	// create execution response struct and return that
	return model.Transaction{}, nil
}
