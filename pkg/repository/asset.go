package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Asset interface {
	Transactable
	Create(model.Asset) (model.Asset, error)
	GetID(id string) (model.Asset, error)
	Update(ID string, updates any) error
}

type asset[T any] struct {
	base[T]
}

func NewAsset(db *sqlx.DB) Asset {
	return &asset[model.Asset]{base[model.Asset]{store: db, table: "asset"}}
}

func (n asset[T]) Create(insert model.Asset) (model.Asset, error) {
	m := model.Asset{}
	rows, err := n.store.NamedQuery(`
		INSERT INTO asset (name) 
		VALUES(:name) 	RETURNING *`, insert)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}

	defer rows.Close()
	return m, err
}
