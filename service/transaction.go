package service

import "github.com/String-xyz/string-api/repository"

type Transaction interface {
	Quote() (int, error)
	Execute() (int, error)
}

type transaction struct {
	repository repository.Transaction
}

func NewTransactor(repo repository.Transaction) Transaction {
	return &transaction{repository: repo}
}

func (t transaction) Quote() (int, error) {
	// create quote struct and return that
	return 0, nil
}

func (t transaction) Execute() (int, error) {
	// create execution response struct and return that
	return 0, nil
}
