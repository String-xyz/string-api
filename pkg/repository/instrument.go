package repository

import (
	"database/sql"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Instrument interface {
	Transactable
	Create(model.Instrument) (model.Instrument, error)
	GetById(id string) (model.Instrument, error)
	GetWallet(addr string) (model.Instrument, error)
	Update(ID string, updates any) error
	GetWalletByUserId(userId string) (model.Instrument, error)
	GetBankByUserId(userId string) (model.Instrument, error)
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
		INSERT INTO instrument (type, status, network, public_key, user_id) 
		VALUES(:type, :status, :network, :public_key, :user_id) 	RETURNING *`, insert)
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

func (i instrument[T]) GetWallet(addr string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE public_key = $1", i.table), addr)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) GetWalletByUserId(userId string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'Crypto Wallet'", i.table), userId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) GetBankByUserId(userId string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'Bank Account'", i.table), userId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
