package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	baserepo "github.com/String-xyz/go-lib/v2/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type Identity interface {
	Create(ctx context.Context, insert model.Identity) (identity model.Identity, err error)
	GetById(ctx context.Context, id string) (model.Identity, error)
	List(ctx context.Context, limit int, offset int) ([]model.Identity, error)
	Update(ctx context.Context, id string, updates any) (identity model.Identity, err error)
	GetByUserId(ctx context.Context, userId string) (identity model.Identity, err error)
	GetByAccountId(ctx context.Context, accountId string) (identity model.Identity, err error)
}

type identity[T any] struct {
	baserepo.Base[T]
}

func NewIdentity(db database.Queryable) Identity {
	return &identity[model.Identity]{baserepo.Base[model.Identity]{Store: db, Table: "identity"}}
}

func (i identity[T]) Create(ctx context.Context, insert model.Identity) (identity model.Identity, err error) {
	query, args, err := i.Named(`
		INSERT INTO identity (userId)
		VALUES(:userId) RETURNING *`, insert)
	if err != nil {
		return identity, libcommon.StringError(err)
	}

	err = i.Store.QueryRowxContext(ctx, query, args...).StructScan(&m)
	if err != nil {
		return identity, libcommon.StringError(err)
	}

	return identity, nil
}

func (i identity[T]) Update(ctx context.Context, id string, updates any) (identity model.Identity, err error) {
	names, keyToUpdate := libcommon.KeysAndValues(updates)
	if len(names) == 0 {
		return identity, libcommon.StringError(errors.New("no updates provided"))
	}

	keyToUpdate["id"] = id

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=:id RETURNING *", i.Table, strings.Join(names, ","))
	namedQuery, args, err := i.Named(query, keyToUpdate)
	if err != nil {
		return identity, libcommon.StringError(err)
	}

	err = i.Store.QueryRowxContext(ctx, namedQuery, args...).StructScan(&identity)
	if err != nil {
		return identity, libcommon.StringError(err)
	}
	return identity, nil
}

func (i identity[T]) GetByUserId(ctx context.Context, userId string) (identity model.Identity, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE userId=$1", i.Table)
	err = i.Store.QueryRowxContext(ctx, query, userId).StructScan(&identity)
	if err != nil {
		return identity, libcommon.StringError(err)
	}
	return identity, nil
}

func (i identity[T]) GetByAccountId(ctx context.Context, accountId string) (identity model.Identity, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE accountId=$1", i.Table)
	err = i.Store.QueryRowxContext(ctx, query, accountId).StructScan(&identity)
	if err != nil {
		return identity, libcommon.StringError(err)
	}
	return identity, nil
}
