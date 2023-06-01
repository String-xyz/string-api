package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	baserepo "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"

	"github.com/String-xyz/string-api/pkg/model"
)

type User interface {
	database.Transactable
	Create(ctx context.Context, user model.User) (model.User, error)
	GetById(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, limit int, offset int) ([]model.User, error)
	Update(ctx context.Context, id string, updates any) (model.User, error)
	GetByType(ctx context.Context, label string) (model.User, error)
	UpdateStatus(ctx context.Context, id string, status string) (model.User, error)
	GetPlatforms(ctx context.Context, id string, limit int, offset int) ([]model.Platform, error)
	GetWithContact(ctx context.Context, id string) (model.UserWithContact, error)
}

type user[T any] struct {
	baserepo.Base[T]
}

func NewUser(db database.Queryable) User {
	return &user[model.User]{baserepo.Base[model.User]{Store: db, Table: "string_user"}}
}

func (u user[T]) Create(ctx context.Context, insert model.User) (model.User, error) {
	m := model.User{}

	query, args, err := u.Named(`
		INSERT INTO string_user (type, status, first_name, middle_name, last_name) 
		VALUES(:type, :status, :first_name, :middle_name, :last_name) RETURNING *`, insert)

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

func (u user[T]) Update(ctx context.Context, id string, updates any) (model.User, error) {
	names, keyToUpdate := libcommon.KeysAndValues(updates)
	var user model.User
	if len(names) == 0 {
		return user, libcommon.StringError(errors.New("no fields to update"))
	}

	// Add the "id" key to the keyToUpdate map
	keyToUpdate["id"] = id

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = :id RETURNING *", u.Table, strings.Join(names, ", "))
	namedQuery, args, err := u.Named(query, keyToUpdate)
	if err != nil {
		return user, libcommon.StringError(err)
	}

	// Use QueryRowxContext to execute the query with the provided context
	err = u.Store.QueryRowxContext(ctx, namedQuery, args...).StructScan(&user)
	if err != nil {
		return user, libcommon.StringError(err)
	}

	return user, nil
}

// update user status
func (u user[T]) UpdateStatus(ctx context.Context, id string, status string) (model.User, error) {
	m := model.User{}

	// Prepare the named statement with sqlx.NamedStmt
	stmt, err := u.Store.PrepareNamedContext(ctx, fmt.Sprintf("UPDATE %s SET status = :status WHERE id = :id RETURNING *", u.Table))
	if err != nil {
		return m, libcommon.StringError(err)
	}
	defer stmt.Close()

	// Execute the prepared named statement with context and parameters
	err = stmt.GetContext(ctx, &m, map[string]interface{}{
		"id":     id,
		"status": status,
	})
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (u user[T]) GetByType(ctx context.Context, label string) (model.User, error) {
	m := model.User{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE type = $1 LIMIT 1", u.Table)
	err := u.Store.GetContext(ctx, &m, query, label)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}

func (u user[T]) GetPlatforms(ctx context.Context, id string, limit int, offset int) (platforms []model.Platform, err error) {
	if limit == 0 {
		limit = 20
	}

	query := `
	SELECT platform.* FROM platform
	LEFT JOIN user_to_platform
	ON platform.id = user_to_platform.platform_id
	WHERE user_to_platform.user_id = $1
	AND platform.deleted_at IS NULL
	LIMIT $2 OFFSET $3`

	err = u.Store.SelectContext(ctx, &platforms, query, id, limit, offset) // Pass id, limit, and offset as separate parameters
	if err != nil && err == sql.ErrNoRows {
		return platforms, serror.NOT_FOUND
	} else if err != nil {
		return platforms, libcommon.StringError(err)
	}
	return platforms, nil
}

func (u user[T]) GetWithContact(ctx context.Context, userId string) (model.UserWithContact, error) {
	m := model.UserWithContact{}
	query := `
	SELECT u.*, c.data as email FROM string_user u
	LEFT JOIN contact c 
	ON u.id = c.user_id
	WHERE u.id = $1 
	AND c.type = 'email' 
	AND c.status = 'validated'
	LIMIT 1
	`
	err := u.Store.GetContext(ctx, &m, query, userId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, libcommon.StringError(err)
	}
	return m, nil
}
