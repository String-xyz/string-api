package service

import (
	"encoding/json"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type Webhook interface {
	Handle(ctx context.Context, data []byte) error
}

type webhook struct{}

func NewWebhook() Webhook {
	return &webhook{}
}

func (w *webhook) Handle(ctx context.Context, data []byte) error {
	event := checkout.WebhookEvent{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return common.StringError(errors.Newf("error unmarshalling webhook event: %v", err))
	}
	return w.processEvent(ctx, event)
}

func (w *webhook) processEvent(ctx context.Context, event checkout.WebhookEvent) error {
	switch event.Type {

	case checkout.AuthorizationApprovedEvent:
		payload := event.Data.(checkout.AuthorizationApproved)
		return w.authorizationApproved(ctx, payload)

	case checkout.AthorizationDeclinedEvent:
		payload := event.Data.(checkout.AuthorizationDeclined)
		return w.authorizationDeclined(ctx, payload)

	case checkout.PaymentApprovedEvent:
		payload := event.Data.(checkout.PaymentApproved)
		return w.paymentApproved(ctx, payload)

	case checkout.PaymentCapturedEvent:
		payload := event.Data.(checkout.PaymentCaptured)
		return w.paymentCaptured(ctx, payload)

	default:
		// only the events above are supported for now.
		return common.StringError(errors.Newf("unhandled event type: %v", event.Type))
	}
}

func (w webhook) authorizationApproved(ctx context.Context, data checkout.AuthorizationApproved) error {
	log.Info().Msgf("authorization approved: %v", data)
	return nil
}

func (w webhook) authorizationDeclined(ctx context.Context, data checkout.AuthorizationDeclined) error {
	log.Info().Msgf("authorization declined: %v", data)
	return nil
}

func (w webhook) paymentApproved(ctx context.Context, data checkout.PaymentApproved) error {
	log.Info().Msgf("payment approved: %v", data)
	return nil
}

func (w webhook) paymentCaptured(ctx context.Context, data checkout.PaymentCaptured) error {
	log.Info().Msgf("payment captured: %v", data)
	return nil
}
