package repository

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type UserToPlatform interface {
	database.Transactable
	Create(ctx context.Context, m model.UserToPlatform) (model.UserToPlatform, error)
	GetById(ctx context.Context, id string) (model.UserToPlatform, error)
	List(ctx context.Context, limit int, offset int) ([]model.UserToPlatform, error)
	ListByUserId(ctx context.Context, userId string, imit int, offset int) ([]model.UserToPlatform, error)
	Update(ctx context.Context, id string, updates any) error
}

type userToPlatform[T any] struct {
	baserepo.Base[T]
}

func NewUserToPlatform(db database.Queryable) UserToPlatform {
	return &userToPlatform[model.UserToPlatform]{baserepo.Base[model.UserToPlatform]{Store: db, Table: "user_to_platform"}}
}

func (u userToPlatform[T]) Create(ctx context.Context, insert model.UserToPlatform) (model.UserToPlatform, error) {
	m := model.UserToPlatform{}
	query := `INSERT INTO user_to_platform (user_id, platform_id) 
		VALUES (:user_id, :platform_id) RETURNING *`

	stmt, err := u.Store.PrepareNamedContext(ctx, query)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &m, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}
