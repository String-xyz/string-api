package repository

import (
	"context"
	"time"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	baserepo "github.com/String-xyz/go-lib/v2/repository"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx/types"
)

type PlatformUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
}

type Platform interface {
	database.Transactable
	Create(ctx context.Context, m model.Platform) (model.Platform, error)
	GetById(ctx context.Context, id string) (model.Platform, error)
	List(ctx context.Context, limit int, offset int) ([]model.Platform, error)
	Update(ctx context.Context, id string, updates any) error
	AssociateUser(ctx context.Context, userId string, platformId string) error
	AssociateContact(ctx context.Context, contactId string, platformId string) error
}

type platform[T any] struct {
	baserepo.Base[T]
}

func NewPlatform(db database.Queryable) Platform {
	return &platform[model.Platform]{baserepo.Base[model.Platform]{Store: db, Table: "platform"}}
}

func (p platform[T]) Create(ctx context.Context, insert model.Platform) (model.Platform, error) {
	m := model.Platform{}

	query, args, err := p.Named(`
		INSERT INTO platform (name, description) 
		VALUES(:name, :description) RETURNING *`, insert)

	if err != nil {
		return m, libcommon.StringError(err)
	}

	// Use QueryRowxContext to execute the query with the provided context
	err = p.Store.QueryRowxContext(ctx, query, args...).StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (p platform[T]) AssociateUser(ctx context.Context, userId string, platformId string) error {
	_, err := p.Store.ExecContext(ctx, `
		INSERT INTO user_to_platform (user_id, platform_id) 
		VALUES($1, $2)`, userId, platformId)

	if err != nil {
		return libcommon.StringError(err)
	}

	return nil
}

func (p platform[T]) AssociateContact(ctx context.Context, contactId string, platformId string) error {
	_, err := p.Store.ExecContext(ctx, `
		INSERT INTO contact_to_platform (contact_id, platform_id) 
		VALUES($1, $2)`, contactId, platformId)
	if err != nil {
		return libcommon.StringError(err)
	}

	return nil
}
