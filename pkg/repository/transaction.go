package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Create(model.Transaction) error
	GetID(id string) (model.Transaction, error)
}

type transaction struct {
	store *sqlx.DB
}

func NewTransaction(store *sqlx.DB) Transaction {
	return &transaction{store}
}

func (t transaction) Create(model.Transaction) error {
	return nil
}

func (t transaction) GetID(id string) (model.Transaction, error) {
	return model.Transaction{}, nil
}
