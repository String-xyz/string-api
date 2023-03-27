package repository

import (
	"encoding/json"
	"fmt"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/string-api/pkg/model"
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
	CreateAny(key string, val any, expire time.Duration) error
	CreateJWTRefresh(key string, val string) (model.AuthStrategy, error)
	GetUserIdFromRefreshToken(key string) (string, error)
	Get(string) (model.AuthStrategy, error)
	GetKeyString(key string) (string, error)
	Delete(key string) error
}

type auth[T any] struct {
	baserepo.Base[T]
	redis database.RedisStore
}

func NewAuth(redis database.RedisStore, db database.Queryable) AuthStrategy {
	return &auth[model.AuthStrategy]{baserepo.Base[model.AuthStrategy]{Store: db, Table: "auth_strategy"}, redis}
}

// Create creates a strategy with user password/email
// Ideally this should be move to PG instead of redis
func (a auth[T]) Create(authType AuthType, m model.AuthStrategy) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Data), 8)
	if err != nil {
		return libcommon.StringError(err)
	}
	strat := &m
	strat.Data = string(hash)
	return a.redis.Set(strat.ContactData, strat, 0)
}

func (a auth[T]) CreateAny(key string, val any, expire time.Duration) error {
	return a.redis.Set(key, val, expire)
}

// CreateJWTRefresh creates and persists a refresh jwt token
func (a auth[T]) CreateJWTRefresh(key string, userId string) (model.AuthStrategy, error) {
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

func (a auth[T]) Get(key string) (model.AuthStrategy, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return model.AuthStrategy{}, libcommon.StringError(err)
	}
	authStrat := model.AuthStrategy{}
	err = json.Unmarshal(m, &authStrat)
	if err != nil {
		return model.AuthStrategy{}, libcommon.StringError(err)
	}

	return authStrat, nil
}

// return the user id from the refresh token or error if token is invalid or expired
func (a auth[T]) GetUserIdFromRefreshToken(refreshToken string) (string, error) {
	authStrat, err := a.Get(refreshToken)

	if err != nil {
		return "", libcommon.StringError(err)
	}
	// assert token has not expired
	if authStrat.ExpiresAt.Before(time.Now()) {
		return "", libcommon.StringError(fmt.Errorf("refresh token expired"))
	}
	// assert token has not been deactivated
	if authStrat.DeactivatedAt != nil {
		return "", libcommon.StringError(fmt.Errorf("refresh token deactivated at %s", authStrat.DeactivatedAt))
	}
	// if all is well, return the user id
	return authStrat.Data, nil
}

func (a auth[T]) GetKeyString(key string) (string, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return "", libcommon.StringError(err)
	}
	return string(m), nil
}

func (a auth[T]) Delete(key string) error {
	return a.redis.Delete(key)
}
