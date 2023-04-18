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
	Create(ctx context.Context, user model.User) (model.User, error)
	GetById(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, limit int, offset int) ([]model.User, error)
	Update(ctx context.Context, id string, updates any) (model.User, error)
	GetByType(ctx context.Context, label string) (model.User, error)
	UpdateStatus(ctx context.Context, id string, status string) (model.User, error)
}

type user[T any] struct {
	baserepo.Base[T]
}

func NewUser(db database.Queryable) User {
	return &user[model.User]{baserepo.Base[model.User]{Store: db, Table: "string_user"}}
}

func (u user[T]) Create(ctx context.Context, insert model.User) (model.User, error) {
	m := model.User{}

	// Prepare the named statement with sqlx.NamedStmt
	stmt, err := u.Store.PrepareNamedContext(ctx, `
		INSERT INTO string_user (type, status, first_name, middle_name, last_name) 
		VALUES(:type, :status, :first_name, :middle_name, :last_name) RETURNING *`)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	defer stmt.Close()

	// Execute the prepared named statement with context
	rows, err := stmt.QueryxContext(ctx, insert)
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

	// TODO: use prepared statement to avoid sql injection
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
func (u user[T]) UpdateStatus(ctx context.Context, id string, status string) (model.User, error) {
	m := model.User{}

	// Prepare the named statement with sqlx.NamedStmt
	stmt, err := u.Store.PrepareNamedContext(ctx, fmt.Sprintf("UPDATE %s SET status = :status WHERE id = :id RETURNING *", u.Table))
	if err != nil {
		return m, libcommon.StringError(err)
	}
	defer stmt.Close()

	// Execute the prepared named statement with context and parameters
	err = stmt.GetContext(ctx, &m, map[string]interface{}{
		"id":     id,
		"status": status,
	})
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (u user[T]) GetByType(ctx context.Context, label string) (model.User, error) {
	m := model.User{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE type = $1 LIMIT 1", u.Table)
	err := u.Store.GetContext(ctx, &m, query, label)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}
