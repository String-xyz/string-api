package service

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Card interface {
	FetchSavedCards(ctx context.Context, userId string, platformId string) (instruments []checkout.CustomerInstrument, err error)
}

type card struct {
	repos repository.Repositories
}

func NewCard(repos repository.Repositories) Card {
	return &card{repos}
}

func (c card) FetchSavedCards(ctx context.Context, userId string, platformId string) (instruments []checkout.CustomerInstrument, err error) {
	_, finish := Span(ctx, "service.card.FetchSavedCards", "userId", userId)
	defer finish()

	contact, err := c.repos.Contact.GetEmailByUserIdAndPlatformId(ctx, userId, platformId)
	if err != nil {
		return nil, libcommon.StringError(err)
	}
	return GetCustomerInstruments(contact.Data)
}
