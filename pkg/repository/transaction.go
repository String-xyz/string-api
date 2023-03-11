package repository

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type Transaction interface {
	database.Transactable
	Create(model.Transaction) (model.Transaction, error)
	GetById(ctx context.Context, id string) (model.Transaction, error)
	Update(ctx context.Context, id string, updates any) error
}

type transaction[T any] struct {
	baserepo.Base[T]
}

func NewTransaction(db database.Queryable) Transaction {
	return &transaction[model.Transaction]{baserepo.Base[model.Transaction]{Store: db, Table: "transaction"}}
}

func (t transaction[T]) Create(insert model.Transaction) (model.Transaction, error) {
	m := model.Transaction{}
	// TODO: Add platform_id once it becomes available
	rows, err := t.Store.NamedQuery(`
		INSERT INTO transaction (status, network_id, device_id, platform_id, ip_address)
		VALUES(:status, :network_id, :device_id, :platform_id, :ip_address) RETURNING id`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.Scan(&m.Id)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}
