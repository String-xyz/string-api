package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Cost interface {
	Create(model.Cost) error
	GetGas(network string) (model.Cost, error)
	GetToken(token string) (model.Cost, error)
}

type cost struct {
	store *sqlx.DB
}

func NewCost(store *sqlx.DB) Cost {
	return &cost{store}
}

func (c cost) Create(model.Cost) error {
	return nil
}

func (c cost) GetGas(network string) (model.Cost, error) {
	return model.Cost{}, nil
}

func (c cost) GetToken(token string) (model.Cost, error) {
	return model.Cost{}, nil
}
