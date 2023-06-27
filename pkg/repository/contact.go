package repository

import (
	"context"
	"database/sql"
	"fmt"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	"github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Contact interface {
	database.Transactable
	Create(ctx context.Context, m model.Contact) (model.Contact, error)
	GetById(ctx context.Context, id string) (model.Contact, error)
	GetByUserId(ctx context.Context, userId string) (model.Contact, error)
	ListByUserId(ctx context.Context, userId string, imit int, offset int) ([]model.Contact, error)
	List(ctx context.Context, limit int, offset int) ([]model.Contact, error)
	Update(ctx context.Context, id string, updates any) error
	GetByData(ctx context.Context, data string) (model.Contact, error)
	GetEmailByUserIdAndPlatformId(ctx context.Context, userId string, platformId string) (model.Contact, error)
	GetByUserIdAndType(ctx context.Context, userId string, _type string) (model.Contact, error)
	GetByUserIdAndStatus(ctx context.Context, userId string, status string) (model.Contact, error)
}

type contact[T any] struct {
	repository.Base[T]
}

func NewContact(db database.Queryable) Contact {
	return &contact[model.Contact]{repository.Base[model.Contact]{Store: db, Table: "contact"}}
}

func (u contact[T]) Create(ctx context.Context, insert model.Contact) (model.Contact, error) {
	m := model.Contact{}

	query, args, err := u.Named(`
		INSERT INTO contact (user_id, data, type, status) 
		VALUES(:user_id, :data, :type, :status) RETURNING *`, insert)

	if err != nil {
		return m, libcommon.StringError(err)
	}

	// Use QueryRowxContext to execute the query with the provided context
	err = u.Store.QueryRowxContext(ctx, query, args...).StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (u contact[T]) GetByData(ctx context.Context, data string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE data = $1", u.Table), data)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}

// TODO: replace references to GetByUserIdAndStatus with the following:
func (u contact[T]) GetEmailByUserIdAndPlatformId(ctx context.Context, userId string, platformId string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.GetContext(ctx, &m, fmt.Sprintf(`
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

func (u contact[T]) GetByUserIdAndType(ctx context.Context, userId string, _type string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND type = $2 LIMIT 1", u.Table), userId, _type)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}

	return m, libcommon.StringError(err)
}

func (u contact[T]) GetByUserIdAndStatus(ctx context.Context, userId, status string) (model.Contact, error) {
	m := model.Contact{}
	err := u.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND status = $2 LIMIT 1", u.Table), userId, status)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}

	return m, libcommon.StringError(err)
}
