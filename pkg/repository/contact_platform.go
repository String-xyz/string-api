package repository

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type ContactPlatform interface {
	Transactable
	Readable
	Create(model.ContactPlatform) (model.ContactPlatform, error)
	GetById(ID string) (model.ContactPlatform, error)
	List(limit int, offset int) ([]model.ContactPlatform, error)
	Update(ID string, updates any) error
}

type contactPlatform[T any] struct {
	base[T]
}

func NewContactPlatform(db *sqlx.DB) ContactPlatform {
	return &contactPlatform[model.ContactPlatform]{base: base[model.ContactPlatform]{store: db, table: "contact_platform"}}
}

func (u contactPlatform[T]) Create(insert model.ContactPlatform) (model.ContactPlatform, error) {
	m := model.ContactPlatform{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO contact_platform (contact_id, platform_id) 
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
