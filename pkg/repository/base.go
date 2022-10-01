package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

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

type Transactable interface {
	MustBegin() Queryable
	Rollback()
	Commit()
	SetTx(t Queryable)
}

type base[T any] struct {
	store Queryable
	db    Queryable
	table string
}

// MustBegin panic is transaction cant start
// the underlying store is set to the transaction
// returned by db.MustBegin()
// You must call rollBack or Commit to return back to a Db state
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
	b.store = b.db
}

func (b *base[T]) Commit() {
	t := b.store.(*sqlx.Tx)
	t.Commit()
	b.store = b.db
}

func (b *base[T]) SetTx(t Queryable) {
	b.store = t
}

func (u base[T]) List(limit int, offset int) (list []T, err error) {
	if limit == 0 {
		limit = 20
	}

	err = u.store.Select(&list, fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", u.table), limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (b base[T]) GetID(ID string) (m T, err error) {
	err = b.store.Get(&m, fmt.Sprintf("SELECT FROM %s WHERE id = $1 deactivated_at = NULL", b.table), ID)
	return m, err
}

func (b base[T]) Update(ID string, updates any) error {
	names, keyToUpdate := keysAndValues(updates)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %s", b.table, strings.Join(names, ","), ID)
	_, err := b.store.NamedExec(query, keyToUpdate)
	return err
}

// keysAndValues is only being used for optional updates
// do not use it for insert or select
func keysAndValues(item interface{}) ([]string, map[string]interface{}) {
	tag := "db"
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	keyNames := make([]string, 0, v.NumField())
	keyValues := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := reflectValue.Field(i).Interface()
		if !isNil(field) {
			t := v.Field(i).Tag.Get(tag) + "=:" + v.Field(i).Tag.Get(tag)
			keyNames = append(keyNames, t)
			keyValues[v.Field(i).Tag.Get(tag)] = field
		}
	}

	return keyNames, keyValues
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
