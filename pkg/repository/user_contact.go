package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type UserContact interface {
	Transactable
	Create(model.Contact) (model.Contact, error)
	GetID(ID string) (model.Contact, error)
	GetUserID(userID string) (model.Contact, error)
	ListUserID(userID string, imit int, offset int) ([]model.Contact, error)
	List(limit int, offset int) ([]model.Contact, error)
	Update(ID string, updates any) error
}

type userContact[T any] struct {
	base[T]
}

func NewUserContact(db *sqlx.DB) UserContact {
	return &userContact[model.Contact]{base: base[model.Contact]{store: db, table: "contact"}}
}

func (u userContact[T]) Create(insert model.Contact) (model.Contact, error) {
	m := model.Contact{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO contact (user_id, data, type, status) 
		VALUES(:user_id, :data, :type, :status) RETURNING *`, insert)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}
	defer rows.Close()
	return m, err
}
