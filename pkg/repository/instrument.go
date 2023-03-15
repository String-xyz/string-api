package repository

import (
	"context"
	"database/sql"
	"fmt"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Instrument interface {
	database.Transactable
	Create(model.Instrument) (model.Instrument, error)
	Update(ctx context.Context, id string, updates any) error
	GetById(ctx context.Context, id string) (model.Instrument, error)
	GetWalletByAddr(addr string) (model.Instrument, error)
	GetCardByFingerprint(fingerprint string) (m model.Instrument, err error)
	GetWalletByUserId(userId string) (model.Instrument, error)
	GetBankByUserId(userId string) (model.Instrument, error)
	WalletAlreadyExists(addr string) (bool, error)
}

type instrument[T any] struct {
	baserepo.Base[T]
}

func NewInstrument(db *sqlx.DB) Instrument {
	return &instrument[model.Instrument]{baserepo.Base[model.Instrument]{Store: db, Table: "instrument"}}
}

func (i instrument[T]) Create(insert model.Instrument) (model.Instrument, error) {
	m := model.Instrument{}
	rows, err := i.Store.NamedQuery(`
		INSERT INTO instrument (type, status, network, public_key, user_id, last_4) 
		VALUES(:type, :status, :network, :public_key, :user_id, :last_4) 	RETURNING *`, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, libcommon.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}

func (i instrument[T]) GetWalletByAddr(addr string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE public_key = $1", i.Table), addr)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) GetCardByFingerprint(fingerprint string) (m model.Instrument, err error) {
	return i.GetWalletByAddr(fingerprint)
}

func (i instrument[T]) GetWalletByUserId(userId string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'Crypto Wallet'", i.Table), userId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) GetBankByUserId(userId string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'Bank Account'", i.Table), userId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) WalletAlreadyExists(addr string) (bool, error) {
	wallet, err := i.GetWalletByAddr(addr)

	if err != nil && errors.Cause(err).Error() != "not found" { // because we are wrapping error and care about its value
		return true, libcommon.StringError(err)
	} else if err == nil && wallet.UserId != "" {
		return true, libcommon.StringError(errors.New("wallet already associated with user"))
	} else if err == nil && wallet.PublicKey == addr {
		return true, libcommon.StringError(errors.New("wallet already exists"))
	}

	return false, nil
}
