package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	GetCardsByUserId(userId string) ([]model.Instrument, error)
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
		INSERT INTO instrument (type, status, network, public_key, user_id, last_4, name) 
		VALUES(:type, :status, :network, :public_key, :user_id, :last_4, :name) 	RETURNING *`, insert)
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
	query := fmt.Sprintf("SELECT * FROM %s WHERE public_key = $1", i.Table)
	err := i.Store.Get(&m, query, addr)

	switch {
	case err == nil:
		return m, nil
	case errors.Is(err, sql.ErrNoRows),
		strings.Contains(errors.Cause(err).Error(), "not found"),
		strings.Contains(errors.Cause(err).Error(), "no rows in result set"):
		return m, serror.NOT_FOUND
	default:
		return m, libcommon.StringError(err)
	}
}

func (i instrument[T]) GetCardByFingerprint(fingerprint string) (m model.Instrument, err error) {
	return i.GetWalletByAddr(fingerprint)
}

func (i instrument[T]) GetWalletByUserId(userId string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'crypto wallet'", i.Table), userId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) GetBankByUserId(userId string) (model.Instrument, error) {
	m := model.Instrument{}
	err := i.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'bank account'", i.Table), userId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (i instrument[T]) WalletAlreadyExists(addr string) (bool, error) {
	wallet, err := i.GetWalletByAddr(addr)

	// not found error means wallet does not exist
	if serror.Is(err, serror.NOT_FOUND) {
		return false, nil
	}

	// throw unknown errors
	if err != nil {
		return false, libcommon.StringError(err)
	}

	if wallet.UserId != "" || wallet.PublicKey == addr {
		return true, nil
	}

	return false, nil
}

func (i instrument[T]) GetCardsByUserId(userId string) ([]model.Instrument, error) {
	var cards []model.Instrument
	err := i.Store.Select(&cards, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = 'credit card' OR type = 'debit card'", i.Table), userId)
	if err != nil && err == sql.ErrNoRows {
		return cards, serror.NOT_FOUND
	} else if err != nil {
		return cards, libcommon.StringError(err)
	}
	return cards, nil
}
