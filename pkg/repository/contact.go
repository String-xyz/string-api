package repository

import (
	"context"
	"database/sql"
	"fmt"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Contact interface {
	database.Transactable
	Create(model.Contact) (model.Contact, error)
	GetById(ctx context.Context, id string) (model.Contact, error)
	GetByUserId(ctx context.Context, userId string) (model.Contact, error)
	ListByUserId(ctx context.Context, userId string, imit int, offset int) ([]model.Contact, error)
	List(ctx context.Context, limit int, offset int) ([]model.Contact, error)
	Update(ctx context.Context, id string, updates any) error
	GetByData(data string) (model.Contact, error)
	GetEmailByUserIdAndPlatformId(userId string, platformId string) (model.Contact, error)
	GetByUserIdAndType(userId string, _type string) (model.Contact, error)
	GetByUserIdAndStatus(userId string, status string) (model.Contact, error)
}

type contact[T any] struct {
	repository.Base[T]
}

func NewContact(db database.Queryable) Contact {
	return &contact[model.Contact]{repository.Base[model.Contact]{Store: db, Table: "contact"}}
}

func (u contact[T]) Create(insert model.Contact) (model.Contact, error) {
	m := model.Contact{}
	rows, err := u.Store.NamedQuery(`
		INSERT INTO contact (user_id, data, type, status) 
		VALUES(:user_id, :data, :type, :status) RETURNING *`, insert)
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

func (u contact[T]) GetByData(data string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE data = $1", u.Table), data)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, nil
}

// TODO: replace references to GetByUserIdAndStatus with the following:
func (u contact[T]) GetEmailByUserIdAndPlatformId(userId string, platformId string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.Get(&m, fmt.Sprintf(`
	SELECT contact.*
		FROM %s
	LEFT JOIN contact_to_platform
		ON contact.id = contact_to_platform.contact_id
	LEFT JOIN platform
		ON contact_to_platform.platform_id = platform.id
	WHERE contact.type = 'email'
		AND contact.user_id = $1
		AND platform.id = $2
	`, u.Table), userId, platformId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}

func (u contact[T]) GetByUserIdAndType(userId string, _type string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = $2 LIMIT 1", u.Table), userId, _type)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}

func (u contact[T]) GetByUserIdAndStatus(userId, status string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND status = $2 LIMIT 1", u.Table), userId, status)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}
