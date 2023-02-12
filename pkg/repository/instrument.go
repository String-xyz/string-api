package repository

import (
	"database/sql"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Instrument interface {
	Transactable
	Create(model.Instrument) (model.Instrument, error)
	Update(ID string, updates any) error
	GetById(id string) (model.Instrument, error)
	GetWalletByAddr(addr string) (model.Instrument, error)
	GetCardByFingerprint(fingerprint string) (m model.Instrument, err error)
	GetWalletByUserId(userId string) (model.Instrument, error)
	GetBankByUserId(userId string) (model.Instrument, error)
	WalletAlreadyExists(addr string) (bool, error)
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
		INSERT INTO instrument (type, status, network, public_key, user_id, last_4) 
		VALUES(:type, :status, :network, :public_key, :user_id, :last_4) 	RETURNING *`, insert)
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

func (i instrument[T]) GetWalletByAddr(addr string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE public_key = $1", i.table), addr)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) GetCardByFingerprint(fingerprint string) (m model.Instrument, err error) {
	m, err = i.GetWalletByAddr(fingerprint)
	if err != nil {
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

func (i instrument[T]) WalletAlreadyExists(addr string) (bool, error) {
	wallet, err := i.GetWalletByAddr(addr)

	if err != nil && errors.Cause(err).Error() != "not found" { // because we are wrapping error and care about its value
		return true, common.StringError(err)
	} else if err == nil && wallet.UserID != "" {
		return true, common.StringError(errors.New("wallet already associated with user"))
	} else if err == nil && wallet.PublicKey == addr {
		return true, common.StringError(errors.New("wallet already exists"))
	}

	return false, nil
}
