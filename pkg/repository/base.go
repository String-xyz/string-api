package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/String-xyz/string-api/pkg/internal/common"
)

var ErrNotFound = errors.New("not found")

type Queryable interface {
	sqlx.Ext
	sqlx.ExecerContext
	sqlx.PreparerContext
	sqlx.QueryerContext
	sqlx.Preparer

	GetContext(context.Context, interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	MustExecContext(context.Context, string, ...interface{}) sql.Result
	PreparexContext(context.Context, string) (*sqlx.Stmt, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Select(interface{}, string, ...interface{}) error
	QueryRow(string, ...interface{}) *sql.Row
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Preparex(string) (*sqlx.Stmt, error)
	NamedExec(string, interface{}) (sql.Result, error)
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	MustExec(string, ...interface{}) sql.Result
	NamedQuery(string, interface{}) (*sqlx.Rows, error)
}

type Readable interface {
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
}

type Transactable interface {
	// MustBegin panic if Tx cant start
	// the underlying store is set to *sqlx.Tx
	// You must call rollBack(), Commit() or Reset() to return back from *sqlx.Tx to *sqlx.DB
	MustBegin() Queryable
	// Rollback rollback the underyling Tx and resets back to  *sqlx.DB from *sqlx.Tx
	Rollback()
	// Commit commits the undelying Tx and resets to back to *sqlx.DB from *sqlx.Tx
	Commit() error
	// SetTx sets the underying store to be sqlx.TX so it can be used for transaction across multiple repos
	SetTx(t Queryable)
	// Reset changes the store back to *sqlx.DB from *sqlx.Tx
	// Useful when there are many repos using the same *sqlx.Tx
	Reset(b ...base[any])
}

type base[T any] struct {
	store Queryable
	db    Queryable
	table string
}

func (b *base[T]) MustBegin() Queryable {
	db := b.store.(*sqlx.DB)
	b.db = db
	t := db.MustBegin()
	b.store = t
	return t
}

func (b *base[T]) Rollback() {
	t := b.store.(*sqlx.Tx)
	t.Rollback()
	b.Reset()
}

func (b *base[T]) Commit() error {
	t := b.store.(*sqlx.Tx)
	err := t.Commit()
	b.Reset()
	if err != nil {
		common.StringError(err)
	}
	return nil
}

func (b *base[T]) SetTx(t Queryable) {
	b.store = t
}

func (b *base[T]) Reset(repos ...base[any]) {
	b.store = b.db
	for _, v := range repos {
		v.Reset()
	}
}

func (b base[T]) List(limit int, offset int) (list []T, err error) {
	if limit == 0 {
		limit = 20
	}

	err = b.store.Select(&list, fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", b.table), limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (b base[T]) GetById(ID string) (m T, err error) {
	err = b.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND 'deactivated_at' IS NOT NULL", b.table), ID)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	}
	return m, nil
}

// Returns the first match of the user's ID
func (b base[T]) GetByUserId(userID string) (m T, err error) {
	err = b.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND 'deactivated_at' IS NOT NULL LIMIT 1", b.table), userID)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	}
	return m, nil
}

func (b base[T]) ListByUserId(userID string, limit int, offset int) ([]T, error) {
	list := []T{}
	if limit == 0 {
		limit = 20
	}
	err := b.store.Select(&list, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 LIMIT $2 OFFSET $3", b.table), userID, limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	if err != nil {
		return list, common.StringError(err)
	}
	return list, nil
}

func (b base[T]) Update(ID string, updates any) error {
	names, keyToUpdate := common.KeysAndValues(updates)
	if len(names) == 0 {
		return common.StringError(errors.New("no fields to update"))
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s'", b.table, strings.Join(names, ", "), ID)
	_, err := b.store.NamedExec(query, keyToUpdate)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func (b base[T]) Select(model interface{}, query string, params ...interface{}) error {
	return b.store.Select(model, query, params)
}

func (b base[T]) Get(model interface{}, query string, params ...interface{}) error {
	return b.store.Get(model, query, params)
}
