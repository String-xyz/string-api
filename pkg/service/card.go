package service

import (
	"context"

	libcommon "github.com/String-xyz/go-lib/v2/common"

	"github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Card interface {
	ListByUserId(ctx context.Context, userId string, platformId string) (instruments []checkout.CardInstrument, err error)
}

type card struct {
	repos repository.Repositories
}

func NewCard(repos repository.Repositories) Card {
	return &card{repos}
}

func (c card) ListByUserId(ctx context.Context, userId string, platformId string) (instruments []checkout.CardInstrument, err error) {
	_, finish := Span(ctx, "service.card.ListByUserId", SpanTag{"platformId": platformId})
	defer finish()

	user, err := c.repos.User.GetById(ctx, userId)
	if err != nil {
		return nil, libcommon.StringError(err)
	}

	client := checkout.New()
	instruments, err = client.Customer.ListInstruments(user.CheckoutId)

	return instruments, libcommon.StringError(err)
}
