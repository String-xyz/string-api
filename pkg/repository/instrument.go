package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Instrument interface {
	Transactable
	Create(model.Instrument) (model.Instrument, error)
	GetID(id string) (model.Instrument, error)
	Update(ID string, updates any) error
}

type instrument[T any] struct {
	base[T]
}

func NewInstrument(db *sqlx.DB) Instrument {
	return &instrument[model.Instrument]{base[model.Instrument]{store: db, table: "instrument"}}
}

func (i instrument[T]) Create(insert model.Instrument) (model.Instrument, error) {
	m := model.Instrument{}
	rows, err := i.store.NamedQuery(`
		INSERT INTO instrument (name) 
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
