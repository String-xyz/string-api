package repository

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type UserPlatform interface {
	Transactable
	Readable
	Create(model.UserPlatform) (model.UserPlatform, error)
	GetById(ID string) (model.UserPlatform, error)
	List(limit int, offset int) ([]model.UserPlatform, error)
	ListByUserId(userID string, imit int, offset int) ([]model.UserPlatform, error)
	Update(ID string, updates any) error
}

type userPlatform[T any] struct {
	base[T]
}

func NewUserPlatform(db *sqlx.DB) UserPlatform {
	return &userPlatform[model.UserPlatform]{base: base[model.UserPlatform]{store: db, table: "user_platform"}}
}

func (u userPlatform[T]) Create(insert model.UserPlatform) (model.UserPlatform, error) {
	m := model.UserPlatform{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO user_platform (user_id, platform_id) 
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
