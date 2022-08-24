package service

import "github.com/String-xyz/string-api/repository"

type Transactor interface {
	Quote() (int, error)
	Execute() (int, error)
}

type transactor struct {
	repository repository.Transaction
}

func NewTransactor(repo repository.Transaction) Transactor {
	return &transactor{repository: repo}
}

func (t transactor) Quote() (int, error) {
	// create quote struct and return that
	return 0, nil
}

func (t transactor) Execute() (int, error) {
	// create execution response struct and return that
	return 0, nil
}
