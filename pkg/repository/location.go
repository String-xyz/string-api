package repository

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Location interface {
	Transactable
	Create(model.Location) (model.Location, error)
	GetById(id string) (model.Location, error)
	Update(id string, updates any) error
}

type location[T any] struct {
	base[T]
}

func NewLocation(db *sqlx.DB) Location {
	return &location[model.Location]{base[model.Location]{store: db, table: "location"}}
}

func (i location[T]) Create(insert model.Location) (model.Location, error) {
	m := model.Location{}
	rows, err := i.store.NamedQuery(`
		INSERT INTO location (name) 
		VALUES(:name) 	RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}
