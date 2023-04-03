package repository

import (
	"context"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
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
	Create(model.Platform) (model.Platform, error)
	GetById(ctx context.Context, id string) (model.Platform, error)
	List(ctx context.Context, limit int, offset int) ([]model.Platform, error)
	Update(ctx context.Context, id string, updates any) error
	AssociateUser(ctx context.Context, userId string, platformId string) error
}

type platform[T any] struct {
	baserepo.Base[T]
}

func NewPlatform(db database.Queryable) Platform {
	return &platform[model.Platform]{baserepo.Base[model.Platform]{Store: db, Table: "platform"}}
}

func (p platform[T]) Create(m model.Platform) (model.Platform, error) {
	plat := model.Platform{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO platform (name, description) 
		VALUES(:name, :description) RETURNING *`, m)

	if err != nil {
		return plat, libcommon.StringError(err)
	}

	for rows.Next() {
		err := rows.StructScan(&plat)
		if err != nil {
			return plat, libcommon.StringError(err)
		}
	}
	defer rows.Close()
	return plat, nil
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
