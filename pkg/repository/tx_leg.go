package repository

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type TxLeg interface {
	database.Transactable
	Create(ctx context.Context, m model.TxLeg) (model.TxLeg, error)
	GetById(ctx context.Context, id string) (model.TxLeg, error)
	Update(ctx context.Context, id string, updates any) error
}

type txLeg[T any] struct {
	baserepo.Base[T]
}

func NewTxLeg(db database.Queryable) TxLeg {
	return &txLeg[model.TxLeg]{baserepo.Base[model.TxLeg]{Store: db, Table: "tx_leg"}}
}

func (t txLeg[T]) Create(ctx context.Context, insert model.TxLeg) (model.TxLeg, error) {
	m := model.TxLeg{}

	query := `
		INSERT INTO tx_leg (timestamp, amount, value, asset_id, user_id, instrument_id) 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING *`

	err := t.Store.QueryRowxContext(ctx, query, insert.Timestamp, insert.Amount, insert.Value, insert.AssetId, insert.UserId, insert.InstrumentId).StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}
