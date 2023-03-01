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
	GetById(id string) (model.Contact, error)
	GetByUserId(userId string) (model.Contact, error)
	ListByUserId(userId string, imit int, offset int) ([]model.Contact, error)
	List(limit int, offset int) ([]model.Contact, error)
	Update(id string, updates any) error
	GetByData(data string) (model.Contact, error)
	GetByUserIdAndPlatformId(userId string, platformId string) (model.Contact, error)
	GetByUserIdAndStatus(userId string, status string) (model.Contact, error)
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

// TODO: replace references to GetByUserIdAndStatus with the following:
func (u contact[T] GetByUserIdAndPlatformId(userId string, platformId string) (model.Contact, error) {
	m := model.Contact{}
	err := u.store.Get(&m, fmt.Sprintf("SELECT contact.*
	FROM contact
	LEFT JOIN contact_platform
	ON contact.id = contact_to_platform.contact_id
	LEFT JOIN platform
	ON contact_to_platform.platform_id = platform.id
	WHERE contact.user_id = $1
	AND platform.id = $2", u.table), userId, platformId)
	if err != nil && err == sql.ErrNoRows {
		return m, ErrNotFound
	}
	return m, common.StringError(err)
}

func (u contact[T]) GetByUserIdAndStatus(userId, status string) (model.Contact, error) {
	m := model.Contact{}
	err := u.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND status = $2 LIMIT 1", u.table), userId, status)
	if err != nil && err == sql.ErrNoRows {
		return m, ErrNotFound
	}
	return m, common.StringError(err)
}
