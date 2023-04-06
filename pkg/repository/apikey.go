package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

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
