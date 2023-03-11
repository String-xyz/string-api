package repository

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
)

type UserToPlatform interface {
	database.Transactable
	Create(model.UserToPlatform) (model.UserToPlatform, error)
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

func (u userToPlatform[T]) Create(insert model.UserToPlatform) (model.UserToPlatform, error) {
	m := model.UserToPlatform{}
	rows, err := u.Store.NamedQuery(`
		INSERT INTO user_to_platform (user_id, platform_id) 
		VALUES(:user_id, :platform_id) RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}
	defer rows.Close()
	return m, nil
}
