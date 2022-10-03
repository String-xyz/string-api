package repository

import (
	"encoding/json"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type EntityType string
type AuthType string

const (
	EntityTypePlatform = EntityType("platform")
	EntityTypeUser     = EntityType("user")
	AuthTypeJWT        = AuthType("jwt")
	AuthTypeEmail      = AuthType("email")
	AuthTypeAPIKey     = AuthType("api_key")
)

type AuthStrategy interface {
	Create(authType AuthType, m model.AuthStrategy) error
	CreateAny(key string, val any, expire time.Duration) error
	CreateAPIKey(ID string, apiKey string) error
	CreateJWTRefresh(key string, val string) error
	Get(string) (model.AuthStrategy, error)
	GetKeyString(key string) (string, error)
}

type auth struct {
	store *sqlx.DB
	redis store.RedisStore
}

func NewAuth(redis store.RedisStore, store *sqlx.DB) AuthStrategy {
	return &auth{redis: redis, store: store}
}

func (a auth) Create(authType AuthType, m model.AuthStrategy) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Data), 8)
	if err != nil {
		return err
	}
	strat := &m
	strat.Data = string(hash)
	return a.redis.Set(strat.ContactData, strat, 0)
}

func (a auth) CreateAny(key string, val any, expire time.Duration) error {
	return a.redis.Set(key, val, expire)
}

// CreateAPIKey creates and persists an API Key for a platform
func (a auth) CreateAPIKey(ID string, key string) error {
	m := model.AuthStrategy{
		ID:         ID,
		CreatedAt:  time.Now(),
		Type:       string(AuthTypeAPIKey),
		EntityType: string(EntityTypePlatform),
		Data:       key,
	}

	return a.redis.Set(m.Data, m, 0)
}

// CreateJWTRefresh creates and persists a refresh jwt token
// TODO: include expire time
func (a auth) CreateJWTRefresh(key string, val string) error {
	m := model.AuthStrategy{
		ID:         key,
		CreatedAt:  time.Now(),
		Type:       string(AuthTypeJWT),
		EntityType: string(EntityTypeUser),
		Data:       val,
	}

	return a.redis.Set(key, m, 0)
}

func (a auth) Get(key string) (model.AuthStrategy, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return model.AuthStrategy{}, err
	}
	authStrat := model.AuthStrategy{}
	err = json.Unmarshal(m, &authStrat)

	return authStrat, err
}

func (a auth) GetKeyString(key string) (string, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return "", err
	}
	return string(m), nil
}
