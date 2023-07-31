package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	strrepo "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Apikey interface {
	database.Transactable
	GetByData(ctx context.Context, data string, keyType string) (model.Apikey, error)
}

type apikey[T any] struct {
	strrepo.Base[T]
}

func NewApikey(db database.Queryable) Apikey {
	return &apikey[model.Apikey]{strrepo.Base[model.Apikey]{Store: db, Table: "apikey"}}
}

func (p apikey[T]) GetByData(ctx context.Context, data string, keyType string) (model.Apikey, error) {
	m := model.Apikey{}
	err := p.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE data = $1 AND type = $2", p.Table), data, keyType)
	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
