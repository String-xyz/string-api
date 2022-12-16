package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type User interface {
	Transactable
	Readable
	Create(model.User) (model.User, error)
	GetById(ID string) (model.User, error)
	List(limit int, offset int) ([]model.User, error)
	Update(ID string, updates any) (model.User, error)
	GetByType(label string) (model.User, error)
}

type user[T any] struct {
	base[T]
}

func NewUser(db *sqlx.DB) User {
	return &user[model.User]{base[model.User]{store: db, table: "string_user"}}
}

func (u user[T]) Create(insert model.User) (model.User, error) {
	m := model.User{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO string_user (type, status, first_name, middle_name, last_name) 
		VALUES(:type, :status, :first_name, :middle_name, :last_name) 	RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	return m, nil
}

func (u user[T]) Update(ID string, updates any) (model.User, error) {
	names, keyToUpdate := common.KeysAndValues(updates)
	var user model.User
	if len(names) == 0 {
		return user, common.StringError(errors.New("no fields to update"))
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s' RETURNING *", u.table, strings.Join(names, ", "), ID)
	rows, err := u.store.NamedQuery(query, keyToUpdate)

	if err != nil {
		return user, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&user)
	}

	if err != nil {
		return user, common.StringError(err)
	}
	return user, err
}

func (u user[T]) GetByType(label string) (model.User, error) {
	m := model.User{}
	err := u.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE type = $1 LIMIT 1", u.table), label)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
