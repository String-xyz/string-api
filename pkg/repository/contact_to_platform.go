package repository

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type ContactToPlatform interface {
	database.Transactable
	Create(ctx context.Context, m model.ContactToPlatform) (model.ContactToPlatform, error)
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

func (u contactToPlatform[T]) Create(ctx context.Context, insert model.ContactToPlatform) (model.ContactToPlatform, error) {
	m := model.ContactToPlatform{}
	row := u.Store.QueryRowxContext(ctx, `
		INSERT INTO contact_to_platform (contact_id, platform_id) 
		VALUES($1, $2) RETURNING *`, insert.ContactId, insert.PlatformId)

	err := row.StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}
