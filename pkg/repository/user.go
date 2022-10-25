package repository

import (
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
	Update(ID string, updates any) error
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
		INSERT INTO string_user (first_name, last_name, type, status) 
		VALUES(:first_name,:last_name, :type, :status) 	RETURNING *`, insert)
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
