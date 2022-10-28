package repository

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Contact interface {
	Transactable
	Readable
	Create(model.Contact) (model.Contact, error)
	GetID(ID string) (model.User, error)
	Update(ID string, updates any) error
}

type contact[T any] struct {
	base[T]
}

func NewContact(db *sqlx.DB) User {
	return &user[model.User]{base[model.User]{store: db, table: "string_user"}}
}

func (c contact[T]) Create(insert model.Contact) (model.Contact, error) {
	m := model.Contact{}
	rows, err := c.store.NamedQuery(`
		INSERT INTO contact (type, user_id, status) 
		VALUES(:type, :user_id, :status) 	RETURNING *`, insert)
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
