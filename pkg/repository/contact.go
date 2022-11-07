package repository

import (
	"database/sql"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Contact interface {
	Transactable
	Readable
	Create(model.Contact) (model.Contact, error)
	GetById(ID string) (model.Contact, error)
	GetByUserId(userID string) (model.Contact, error)
	ListByUserId(userID string, imit int, offset int) ([]model.Contact, error)
	List(limit int, offset int) ([]model.Contact, error)
	Update(ID string, updates any) error
	GetByData(data string) (model.Contact, error)
}

type contact[T any] struct {
	base[T]
}

func NewContact(db *sqlx.DB) Contact {
	return &contact[model.Contact]{base: base[model.Contact]{store: db, table: "contact"}}
}

func (u contact[T]) Create(insert model.Contact) (model.Contact, error) {
	m := model.Contact{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO contact (user_id, data, type, status) 
		VALUES(:user_id, :data, :type, :status) RETURNING *`, insert)
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

func (u contact[T]) GetByData(data string) (model.Contact, error) {
	m := model.Contact{}
	err := u.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE data = $1", u.table), data)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	}
	return m, nil
}
