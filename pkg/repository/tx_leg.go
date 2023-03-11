package repository

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type TxLeg interface {
	database.Transactable
	Create(model.TxLeg) (model.TxLeg, error)
	GetById(ctx context.Context, id string) (model.TxLeg, error)
	Update(ctx context.Context, id string, updates any) error
}

type txLeg[T any] struct {
	baserepo.Base[T]
}

func NewTxLeg(db database.Queryable) TxLeg {
	return &txLeg[model.TxLeg]{baserepo.Base[model.TxLeg]{Store: db, Table: "tx_leg"}}
}

func (t txLeg[T]) Create(insert model.TxLeg) (model.TxLeg, error) {
	m := model.TxLeg{}
	rows, err := t.Store.NamedQuery(`
		INSERT INTO tx_leg (timestamp, amount, value, asset_id, user_id, instrument_id) 
		VALUES(:timestamp, :amount, :value, :asset_id, :user_id, :instrument_id) 	RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	return m, err
}
