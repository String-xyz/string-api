package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type EntityType = model.EntityType
type AuthType = model.AuthType

const (
	EntityTypePlatform = EntityType("platform")
	EntityTypeUser     = EntityType("user")
	AuthTypeJWT        = AuthType("jwt")
	AuthTypeEmail      = AuthType("email")
	AuthTypePK         = AuthType("privateKey")
	AuthTypeOTP        = AuthType("otp")
	AuthTypeAPIKey     = AuthType("apiKey")
)

type AuthStrategy interface {
	Create(authType AuthType, m model.AuthStrategy) error
	CreateAny(key string, val any, expire time.Duration) error
	CreateAPIKey(entityId string, authType AuthType, apiKey string, persistOnly bool) (model.AuthStrategy, error)
	CreateJWTRefresh(key string, val string) (model.AuthStrategy, error)
	GetUserIdFromRefreshToken(key string) (string, error)
	Get(string) (model.AuthStrategy, error)
	GetKeyString(key string) (string, error)
	List(limit, offset int) ([]model.AuthStrategy, error)
	ListByStatus(limit, offset int, status string) ([]model.AuthStrategy, error)
	UpdateStatus(Id, status string) (model.AuthStrategy, error)
	Delete(key string) error
}

type auth struct {
	store *sqlx.DB
	redis store.RedisStore
}

func NewAuth(redis store.RedisStore, store *sqlx.DB) AuthStrategy {
	return &auth{redis: redis, store: store}
}

// Create creates a strategy with user password/email
// Ideally this should be move to PG instead of redis
func (a auth) Create(authType AuthType, m model.AuthStrategy) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Data), 8)
	if err != nil {
		return common.StringError(err)
	}
	strat := &m
	strat.Data = string(hash)
	return a.redis.Set(strat.ContactData, strat, 0)
}

func (a auth) CreateAny(key string, val any, expire time.Duration) error {
	return a.redis.Set(key, val, expire)
}

// CreateAPIKey creates and persists an API Key for a platform
func (a auth) CreateAPIKey(entityId string, authType AuthType, key string, persistOnly bool) (model.AuthStrategy, error) {
	// only insert to postgres and skip redis cache
	if persistOnly {
		rows, err := a.store.Queryx("INSERT INTO auth_strategy(type,data) VALUES($1, $2) RETURNING *", authType, key)
		if err == nil {
			m := model.AuthStrategy{}
			var scanErr error
			for rows.Next() {
				scanErr = rows.StructScan(&m)
			}
			return m, scanErr
		}
		return model.AuthStrategy{}, err
	}

	m := model.AuthStrategy{
		EntityId:   entityId,
		CreatedAt:  time.Now(),
		Type:       string(authType),
		EntityType: string(EntityTypePlatform),
		Data:       key,
	}

	return m, a.redis.Set(key, m, 0)
}

// CreateJWTRefresh creates and persists a refresh jwt token
func (a auth) CreateJWTRefresh(key string, userId string) (model.AuthStrategy, error) {
	expireAt := time.Hour * 24 * 7 // 7 days expiration
	m := model.AuthStrategy{
		Id:         key,
		CreatedAt:  time.Now(),
		Type:       string(AuthTypeJWT),
		EntityType: string(EntityTypeUser),
		Data:       userId,
		ExpiresAt:  time.Now().Add(expireAt),
	}

	return m, a.redis.Set(key, m, expireAt)
}

func (a auth) Get(key string) (model.AuthStrategy, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return model.AuthStrategy{}, common.StringError(err)
	}
	authStrat := model.AuthStrategy{}
	err = json.Unmarshal(m, &authStrat)
	if err != nil {
		return model.AuthStrategy{}, common.StringError(err)
	}

	return authStrat, nil
}

// return the user id from the refresh token or error if token is invalid or expired
func (a auth) GetUserIdFromRefreshToken(refreshToken string) (string, error) {
	authStrat, err := a.Get(refreshToken)

	if err != nil {
		return "", common.StringError(err)
	}
	// assert token has not expired
	if authStrat.ExpiresAt.Before(time.Now()) {
		return "", common.StringError(fmt.Errorf("refresh token expired"))
	}
	// assert token has not been deactivated
	if authStrat.DeactivatedAt != nil {
		return "", common.StringError(fmt.Errorf("refresh token deactivated at %s", authStrat.DeactivatedAt))
	}
	// if all is well, return the user id
	return authStrat.Data, nil
}

func (a auth) GetKeyString(key string) (string, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return "", common.StringError(err)
	}
	return string(m), nil
}

// List all the available auth_keys on the postgres db
func (a auth) List(limit, offset int) ([]model.AuthStrategy, error) {
	list := []model.AuthStrategy{}
	err := a.store.Select(&list, "SELECT * FROM auth_strategy LIMIT $1 OFFSET $2", limit, offset)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

// ListByStatus lists all auth_keys with a given status on the postgres db
func (a auth) ListByStatus(limit, offset int, status string) ([]model.AuthStrategy, error) {
	list := []model.AuthStrategy{}
	err := a.store.Select(&list, "SELECT * FROM auth_strategy WHERE status = $1 LIMIT $2 OFFSET $3", status, limit, offset)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

// UpdateStatus updates the status on postgres db and returns the updated row
func (a auth) UpdateStatus(Id, status string) (model.AuthStrategy, error) {
	row := a.store.QueryRowx("UPDATE auth_strategy SET status = $2 WHERE id = $1 RETURNING *", Id, status)
	m := model.AuthStrategy{}
	err := row.StructScan(&m)
	return m, err
}

func (a auth) Delete(key string) error {
	return a.redis.Delete(key)
}
