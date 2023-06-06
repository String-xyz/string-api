package service

import (
	"context"
	"time"

	libcommon "github.com/String-xyz/go-lib/v2/common"

	"github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Card interface {
	FetchSavedCards(ctx context.Context, userId string, platformId string) (cards []model.CardResponse, err error)
}

type card struct {
	repos repository.Repositories
}

func NewCard(repos repository.Repositories) Card {
	return &card{repos}
}

func (c card) FetchSavedCards(ctx context.Context, userId string, platformId string) (cards []model.CardResponse, err error) {
	_, finish := Span(ctx, "service.card.FetchSavedCards", SpanTag{"platformId": platformId})
	defer finish()

	user, err := c.repos.User.GetById(ctx, userId)
	if err != nil {
		return nil, libcommon.StringError(err)
	}

	client := checkout.New()

	instruments, err := client.Customer.ListInstruments(user.CheckoutId)
	if err != nil {
		return nil, libcommon.StringError(err)
	}

	for _, instrument := range instruments {
		card := instrument.GetCardInstrumentResponse

		now := time.Now()
		isCardExpired := card.ExpiryYear < now.Year() || (card.ExpiryYear == now.Year() && card.ExpiryMonth < int(now.Month()))

		cards = append(cards, model.CardResponse{
			Type:        card.Type,
			Id:          card.Id,
			Scheme:      card.Scheme,
			Last4:       card.Last4,
			ExpiryMonth: card.ExpiryMonth,
			ExpiryYear:  card.ExpiryYear,
			Expired:     isCardExpired,
			CardType:    card.CardType,
		})
	}

	return cards, nil
}
