package stubs

import (
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type AuthStrategyRepo struct {
}

func (AuthStrategyRepo) Create(authType repository.AuthType, m model.AuthStrategy) error {
	return nil
}

func (AuthStrategyRepo) CreateAPIKey(ID string, apiKey string) error {
	return nil
}

func (AuthStrategyRepo) CreateJWTRefresh(ID string, token string) error {
	return nil
}
func (AuthStrategyRepo) Get(string) (model.AuthStrategy, error) {
	return model.AuthStrategy{}, nil
}

func (AuthStrategyRepo) CreateAny(key string, val any, expire time.Duration) error {
	return nil
}

func (AuthStrategyRepo) GetKeyString(key string) (string, error) {
	return "", nil
}
