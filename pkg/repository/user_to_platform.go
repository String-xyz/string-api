package repository

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type UserToPlatform interface {
	Transactable
	Readable
	Create(model.UserToPlatform) (model.UserToPlatform, error)
	GetById(ID string) (model.UserToPlatform, error)
	List(limit int, offset int) ([]model.UserToPlatform, error)
	ListByUserId(userID string, imit int, offset int) ([]model.UserToPlatform, error)
	Update(ID string, updates any) error
}

type userToPlatform[T any] struct {
	base[T]
}

func NewUserToPlatform(db *sqlx.DB) UserToPlatform {
	return &userToPlatform[model.UserToPlatform]{base: base[model.UserToPlatform]{store: db, table: "user_to_platform"}}
}

func (u userToPlatform[T]) Create(insert model.UserToPlatform) (model.UserToPlatform, error) {
	m := model.UserToPlatform{}
	rows, err := u.store.NamedQuery(`
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
