package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Transactable
	Create(model.Transaction) (model.Transaction, error)
	GetID(id string) (model.Transaction, error)
	Update(ID string, updates any) error
}

type transaction[T any] struct {
	base[T]
}

func NewTransaction(db *sqlx.DB) Transaction {
	return &transaction[model.Transaction]{base[model.Transaction]{store: db, table: "transaction"}}
}

func (t transaction[T]) Create(insert model.Transaction) (model.Transaction, error) {
	m := model.Transaction{}
	rows, err := t.store.NamedQuery(`
		INSERT INTO transaction (type, status) 
		VALUES(:type, :status) 	RETURNING *`, insert)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}

	defer rows.Close()
	return m, err
}
