package repository

import (
	"context"

	commonlib "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type ContactToPlatform interface {
	database.Transactable
	Create(model.ContactToPlatform) (model.ContactToPlatform, error)
	GetById(ctx context.Context, id string) (model.ContactToPlatform, error)
	List(ctx context.Context, limit int, offset int) ([]model.ContactToPlatform, error)
	Update(ctx context.Context, id string, updates any) error
}

type contactToPlatform[T any] struct {
	repository.Base[T]
}

func NewContactPlatform(db database.Queryable) ContactToPlatform {
	return &contactToPlatform[model.ContactToPlatform]{repository.Base[model.ContactToPlatform]{Store: db, Table: "contact_to_platform"}}
}

func (u contactToPlatform[T]) Create(insert model.ContactToPlatform) (model.ContactToPlatform, error) {
	m := model.ContactToPlatform{}
	rows, err := u.Store.NamedQuery(`
		INSERT INTO contact_to_platform (contact_id, platform_id) 
		VALUES(:contact_id, :platform_id) RETURNING *`, insert)
	if err != nil {
		return m, commonlib.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, commonlib.StringError(err)
		}
	}
	defer rows.Close()
	return m, nil
}
