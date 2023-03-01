package repository

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type ContactToPlatform interface {
	Transactable
	Readable
	Create(model.ContactToPlatform) (model.ContactToPlatform, error)
	GetById(id string) (model.ContactToPlatform, error)
	List(limit int, offset int) ([]model.ContactToPlatform, error)
	Update(id string, updates any) error
}

type contactToPlatform[T any] struct {
	base[T]
}

func NewContactPlatform(db *sqlx.DB) ContactToPlatform {
	return &contactToPlatform[model.ContactToPlatform]{base: base[model.ContactToPlatform]{store: db, table: "contact_to_platform"}}
}

func (u contactToPlatform[T]) Create(insert model.ContactToPlatform) (model.ContactToPlatform, error) {
	m := model.ContactToPlatform{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO contact_to_platform (contact_id, platform_id) 
		VALUES(:contact_id, :platform_id) RETURNING *`, insert)
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
