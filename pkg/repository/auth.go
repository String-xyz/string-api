package repository

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/google/uuid"
)

type Auth interface {
	CreateStrategy(model.AuthStrategy) error
	CreateAPIKey(apiKey string) error
	CreateJWTRefresh(token string) error
	GetStrategy(string) (model.AuthStrategy, error)
}

type auth struct {
	redis store.RedisStore
}

func NewAuth(redis store.RedisStore) Auth {
	return &auth{redis}
}

func (a auth) CreateStrategy(m model.AuthStrategy) error {
	return a.redis.Set(m.Token, m)
}

func (a auth) CreateAPIKey(key string) error {
	bs := sha256.Sum256([]byte(key))
	hash := string(bs[:])
	m := model.AuthStrategy{
		ID:         uuid.NewString(),
		CreatedAt:  time.Now(),
		AuthType:   "API_KEY",
		EntityType: "PLATFORM",
		Token:      hash,
	}

	return a.redis.Set(m.Token, m)
}

// CreateJWTRefresh creates a refresh jwt token
func (a auth) CreateJWTRefresh(token string) error {
	bs := sha256.Sum256([]byte(token))
	hash := string(bs[:])
	m := model.AuthStrategy{
		ID:         uuid.NewString(),
		CreatedAt:  time.Now(),
		AuthType:   "JWT_REFRESH",
		EntityType: "USER",
		Token:      hash,
	}

	return a.redis.Set(m.Token, m)
}

// GetStrategy will hash the key and attemp a look up on redis
func (a auth) GetStrategy(key string) (model.AuthStrategy, error) {
	bs := sha256.Sum256([]byte(key))
	m, err := a.redis.Get(string(bs[:]))
	if err != nil {
		return model.AuthStrategy{}, err
	}
	authStrat := model.AuthStrategy{}
	err = json.Unmarshal(m, &authStrat)

	return authStrat, err
}
