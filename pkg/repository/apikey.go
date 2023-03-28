package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type ApikeyUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Data          *string    `json:"data" db:"data"`
	Description   *string    `json:"description" db:"description"`
	CreatedBy     *string    `json:"createdBy" db:"created_by"`
	PlatformID    *string    `json:"platformId" db:"platform_id"`
}

type Apikey interface {
	database.Transactable
	GetByData(ctx context.Context, data string) (model.Apikey, error)
}

type apikey[T any] struct {
	strrepo.Base[T]
}

func NewApikey(db database.Queryable) Apikey {
	return &apikey[model.Apikey]{strrepo.Base[model.Apikey]{Store: db, Table: "apikey"}}
}

func (p apikey[T]) GetByData(ctx context.Context, data string) (model.Apikey, error) {
	m := model.Apikey{}
	err := p.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE data = $1", p.Table), data)
	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
