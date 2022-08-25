package service

import (
	"github.com/String-xyz/string-api/model"
	"github.com/String-xyz/string-api/repository"
)

type Transaction interface {
	Quote() (model.SignedQuote, error)
	Execute() (model.TransactionResponse, error)
}

type transaction struct {
	repository repository.Transaction
}

func NewTransaction(repo repository.Transaction) Transaction {
	return &transaction{repository: repo}
}

func (t transaction) Quote() (model.SignedQuote, error) {
	// create quote struct and return that
	return model.SignedQuote{}, nil
}

func (t transaction) Execute() (model.TransactionResponse, error) {

	// create execution response struct and return that
	return model.TransactionResponse{}, nil
}
