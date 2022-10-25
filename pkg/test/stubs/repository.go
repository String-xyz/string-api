package stubs

import (
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

// type Asset struct {

// }

// var avax = model.Asset{
// 	ID:          "1",
// 	CreatedAt:   time.Now(),
// 	UpdatedAt:   time.Now(),
// 	Name:        "AVAX",
// 	Description: "Avalanche",
// 	Decimals:    18,
// 	IsCrypto:    true,
// 	NetworkID:   sql.NullString{},
// 	ValueOracle: sql.NullString{String: "avalanche-2", Valid: true},
// }

// var usd = model.Asset{
// 	ID:          "2",
// 	CreatedAt:   time.Now(),
// 	UpdatedAt:   time.Now(),
// 	Name:        "USD",
// 	Description: "United States Dollar",
// 	Decimals:    6,
// 	IsCrypto:    false,
// 	NetworkID:   sql.NullString{},
// 	ValueOracle: sql.NullString{},
// }

// func (Asset) Create(model.Asset) (model.Asset, error) {
// 	return model.Asset{}, nil
// }

// func (Asset) GetById(id string) (model.Asset, error) {
// 	if id == "1" {
// 		return avax, nil
// 	}
// 	if id == "2" {
// 		return usd, nil
// 	}
// 	return model.Asset{}, nil
// }

// func (Asset) GetName(name string) (model.Asset, error) {
// 	if name == "AVAX" {
// 		return avax, nil
// 	}
// 	if name == "USD" {
// 		return usd, nil
// 	}
// 	return model.Asset{}, nil
// }

// func (Asset) Update(ID string, updates any) error {
// 	return nil
// }

type AuthStrategyRepo struct {
}

func (AuthStrategyRepo) Create(authType repository.AuthType, m model.AuthStrategy) error {
	return nil
}

func (AuthStrategyRepo) CreateAPIKey(entityID string, authType model.AuthType, apiKey string) error {
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
