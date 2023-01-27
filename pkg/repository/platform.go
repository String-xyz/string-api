package repository

import (
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type PlaformUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
}

type Platform interface {
	Transactable
	Create(model.Platform) (model.Platform, error)
	GetById(ID string) (model.Platform, error)
	List(limit int, offset int) ([]model.Platform, error)
	Update(ID string, updates any) error
}

type platform[T any] struct {
	base[T]
}

func NewPlatform(db *sqlx.DB) Platform {
	return &platform[model.Platform]{base: base[model.Platform]{store: db, table: "platform"}}
}

func (p platform[T]) Create(m model.Platform) (model.Platform, error) {
	plat := model.Platform{}
	rows, err := p.store.NamedQuery(`
		INSERT INTO platform (name, description) 
		VALUES(:name, :description) RETURNING *`, m)

	if err != nil {
		return plat, common.StringError(err)
	}

	for rows.Next() {
		err := rows.StructScan(&plat)
		if err != nil {
			return plat, common.StringError(err)
		}
	}
	defer rows.Close()
	return plat, nil
}

