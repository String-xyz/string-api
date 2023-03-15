package repository

import (
	"context"
	"time"

	libCommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx/types"
)

type PlaformUpdates struct {
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
		return plat, libCommon.StringError(err)
	}

	for rows.Next() {
		err := rows.StructScan(&plat)
		if err != nil {
			return plat, libCommon.StringError(err)
		}
	}
	defer rows.Close()
	return plat, nil
}
