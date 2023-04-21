package repository

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Location interface {
	database.Transactable
	Create(ctx context.Context, m model.Location) (model.Location, error)
	GetById(ctx context.Context, id string) (model.Location, error)
	Update(ctx context.Context, id string, updates any) error
}

type location[T any] struct {
	baserepo.Base[T]
}

func NewLocation(db *sqlx.DB) Location {
	return &location[model.Location]{baserepo.Base[model.Location]{Store: db, Table: "location"}}
}

func (i location[T]) Create(ctx context.Context, insert model.Location) (model.Location, error) {
	// TODO: Use ctx query

	m := model.Location{}
	rows, err := i.Store.NamedQuery(`
		INSERT INTO location (name) 
		VALUES(:name) 	RETURNING *`, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, libcommon.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}
