package repository

import (
	"database/sql"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Asset interface {
	Transactable
	Create(model.Asset) (model.Asset, error)
	GetID(id string) (model.Asset, error)
	GetName(name string) (model.Asset, error)
	Update(ID string, updates any) error
}

type asset[T any] struct {
	base[T]
}

func NewAsset(db *sqlx.DB) Asset {
	return &asset[model.Asset]{base[model.Asset]{store: db, table: "asset"}}
}

func (a asset[T]) Create(insert model.Asset) (model.Asset, error) {
	m := model.Asset{}
	rows, err := a.store.NamedQuery(`
		INSERT INTO asset (name) 
		VALUES(:name) 	RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}

	defer rows.Close()
	return m, err
}

func (a asset[T]) GetName(name string) (model.Asset, error) {
	m := model.Asset{}
	err := a.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE name = $1", a.table), name)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	}
	return m, nil
}
