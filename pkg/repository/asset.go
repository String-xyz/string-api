package repository

import (
	"context"
	"database/sql"
	"fmt"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Asset interface {
	database.Transactable
	Create(ctx context.Context, m model.Asset) (model.Asset, error)
	GetById(ctx context.Context, id string) (model.Asset, error)
	GetByName(ctx context.Context, name string) (model.Asset, error)
	Update(ctx context.Context, Id string, updates any) error
}

type asset[T any] struct {
	baserepo.Base[T]
}

func NewAsset(db database.Queryable) Asset {
	return &asset[model.Asset]{baserepo.Base[model.Asset]{Store: db, Table: "asset"}}
}

func (a asset[T]) Create(ctx context.Context, insert model.Asset) (model.Asset, error) {
	m := model.Asset{}

	query, args, err := a.Named(`
		INSERT INTO asset (name, description, decimals, is_crypto, network_id, value_oracle, value_oracle_2) 
		VALUES(:name, :description, :decimals, :is_crypto, :network_id, :value_oracle, :value_oracle_2) RETURNING *`, insert)

	if err != nil {
		return m, libcommon.StringError(err)
	}

	// Use QueryRowxContext to execute the query with the provided context
	err = a.Store.QueryRowxContext(ctx, query, args...).StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (a asset[T]) GetByName(ctx context.Context, name string) (model.Asset, error) {
	m := model.Asset{}
	err := a.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE name = $1", a.Table), name)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, nil
}
