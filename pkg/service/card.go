package service

import (
	"context"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Card interface {
	FetchSavedCards(ctx context.Context, userId string) ([]model.Instrument, error)
}

type card struct {
	repos repository.Repositories
}

func NewCard(repos repository.Repositories) Card {
	return &card{repos}
}

func (c card) FetchSavedCards(ctx context.Context, userId string) ([]model.Instrument, error) {
	return c.repos.Instrument.GetCardsByUserId(userId)
}
