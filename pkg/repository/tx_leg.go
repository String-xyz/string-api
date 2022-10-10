package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type TxLeg interface {
	Transactable
	Create(model.TxLeg) (model.TxLeg, error)
	GetID(id string) (model.TxLeg, error)
	Update(ID string, updates any) error
}

type txLeg[T any] struct {
	base[T]
}

func NewTxLeg(db *sqlx.DB) TxLeg {
	return &txLeg[model.TxLeg]{base[model.TxLeg]{store: db, table: "tx_leg"}}
}

func (t txLeg[T]) Create(insert model.TxLeg) (model.TxLeg, error) {
	m := model.TxLeg{}
	rows, err := t.store.NamedQuery(`
		INSERT INTO tx_leg (instrument_id) 
		VALUES(:instrument_id) 	RETURNING *`, insert)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}

	defer rows.Close()
	return m, err
}
