package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Asset interface {
	database.Transactable
	Create(model.Asset) (model.Asset, error)
	GetById(ctx context.Context, id string) (model.Asset, error)
	GetByName(name string) (model.Asset, error)
	Update(ctx context.Context, Id string, updates any) error
}

type asset[T any] struct {
	baserepo.Base[T]
}

func NewAsset(db database.Queryable) Asset {
	return &asset[model.Asset]{baserepo.Base[model.Asset]{Store: db, Table: "asset"}}
}

func (a asset[T]) Create(insert model.Asset) (model.Asset, error) {
	m := model.Asset{}
	rows, err := a.Store.NamedQuery(`
		INSERT INTO asset (name, description, decimals, is_crypto, network_id, value_oracle) 
		VALUES(:name, :description, :decimals, :is_crypto, :network_id, :value_oracle) 	RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}

	defer rows.Close()
	return m, err
}

func (a asset[T]) GetByName(name string) (model.Asset, error) {
	m := model.Asset{}
	err := a.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE name = $1", a.Table), name)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, nil
}
