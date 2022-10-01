package repository

import (
	"crypto/sha256"
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
	CreateAPIKey(ID string, apiKey string) error
	CreateJWTRefresh(ID string, token string) error
	Get(string) (model.AuthStrategy, error)
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
	return a.redis.Set(m.ContactData, strat)
}

// CreateAPIKey creates and persists an API Key for a platform
func (a auth) CreateAPIKey(ID string, key string) error {
	bs := sha256.Sum256([]byte(key))
	hash := string(bs[:])
	m := model.AuthStrategy{
		ID:         ID,
		CreatedAt:  time.Now(),
		Type:       string(AuthTypeAPIKey),
		EntityType: string(EntityTypePlatform),
		Data:       hash,
	}

	return a.redis.Set(m.Data, m)
}

// CreateJWTRefresh creates and persists a refresh jwt token
func (a auth) CreateJWTRefresh(ID string, token string) error {
	bs := sha256.Sum256([]byte(token))
	hash := string(bs[:])
	m := model.AuthStrategy{
		ID:         ID,
		CreatedAt:  time.Now(),
		Type:       string(AuthTypeJWT),
		EntityType: string(EntityTypeUser),
		Data:       hash,
	}

	return a.redis.Set(m.Data, m)
}

// Get will hash the key and attemp a look up on redis using the key(JWT refresh token | API key)
func (a auth) Get(key string) (model.AuthStrategy, error) {
	bs := sha256.Sum256([]byte(key))
	m, err := a.redis.Get(string(bs[:]))
	if err != nil {
		return model.AuthStrategy{}, err
	}
	authStrat := model.AuthStrategy{}
	err = json.Unmarshal(m, &authStrat)

	return authStrat, err
}
