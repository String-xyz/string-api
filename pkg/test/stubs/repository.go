package stubs

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type AuthStrategyRepoStub struct {
}

func (a AuthStrategyRepoStub) Create(authType repository.AuthType, m model.AuthStrategy) error {
	return nil
}
func (a AuthStrategyRepoStub) CreateAPIKey(ID string, apiKey string) error {
	return nil
}
func (a AuthStrategyRepoStub) CreateJWTRefresh(ID string, token string) error {
	return nil
}
func (a AuthStrategyRepoStub) Get(string) (model.AuthStrategy, error) {
	return model.AuthStrategy{}, nil
}
