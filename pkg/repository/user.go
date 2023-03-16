package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type User interface {
	database.Transactable
	Create(model.User) (model.User, error)
	GetById(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, limit int, offset int) ([]model.User, error)
	Update(ctx context.Context, id string, updates any) (model.User, error)
	GetByType(label string) (model.User, error)
	UpdateStatus(id string, status string) (model.User, error)
}

type user[T any] struct {
	baserepo.Base[T]
}

func NewUser(db database.Queryable) User {
	return &user[model.User]{baserepo.Base[model.User]{Store: db, Table: "string_user"}}
}

func (u user[T]) Create(insert model.User) (model.User, error) {
	m := model.User{}
	rows, err := u.Store.NamedQuery(`
		INSERT INTO string_user (type, status, first_name, middle_name, last_name) 
		VALUES(:type, :status, :first_name, :middle_name, :last_name) 	RETURNING *`, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, libcommon.StringError(err)
		}
	}

	return m, nil
}

func (u user[T]) Update(ctx context.Context, id string, updates any) (model.User, error) {
	names, keyToUpdate := libcommon.KeysAndValues(updates)
	var user model.User
	if len(names) == 0 {
		return user, libcommon.StringError(errors.New("no fields to update"))
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s' RETURNING *", u.Table, strings.Join(names, ", "), id)
	rows, err := u.Store.NamedQuery(query, keyToUpdate)

	if err != nil {
		return user, libcommon.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&user)
	}

	if err != nil {
		return user, libcommon.StringError(err)
	}
	return user, err
}

// update user status
func (u user[T]) UpdateStatus(id string, status string) (model.User, error) {
	m := model.User{}
	err := u.Store.Get(&m, fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 RETURNING *", u.Table), status, id)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (u user[T]) GetByType(label string) (model.User, error) {
	m := model.User{}
	err := u.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE type = $1 LIMIT 1", u.Table), label)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}
