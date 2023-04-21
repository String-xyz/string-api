package repository

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	database.Transactable
	Create(ctx context.Context, m model.Transaction) (model.Transaction, error)
	GetById(ctx context.Context, id string) (model.Transaction, error)
	Update(ctx context.Context, id string, updates any) error
}

type transaction[T any] struct {
	baserepo.Base[T]
}

func NewTransaction(db database.Queryable) Transaction {
	return &transaction[model.Transaction]{baserepo.Base[model.Transaction]{Store: db, Table: "transaction"}}
}

func (t transaction[T]) Create(ctx context.Context, insert model.Transaction) (model.Transaction, error) {
	m := model.Transaction{}

	query := `
		INSERT INTO transaction (status, network_id, device_id, platform_id, ip_address)
		VALUES (:status, :network_id, :device_id, :platform_id, :ip_address) RETURNING id`

	// Prepare the named query
	query, args, err := sqlx.Named(query, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	// Replace the named parameters with placeholders
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	// Rebind the query to adapt placeholders to the specific SQL database
	query = t.Store.Rebind(query)

	// Use QueryRowxContext to execute the query with the provided context
	err = t.Store.QueryRowxContext(ctx, query, args...).Scan(&m.Id)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}
