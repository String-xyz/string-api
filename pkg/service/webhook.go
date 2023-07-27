package service

import (
	"encoding/json"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"

	"github.com/String-xyz/string-api/pkg/internal/checkout"
)

type WebhookType string

const (
	WebhookTypePersona  WebhookType = "persona"
	WebhookTypeCheckout WebhookType = "checkout"
)

type Webhook interface {
	Handle(ctx context.Context, data []byte, webhook WebhookType) error
}

type webhook struct {
	person personaWebhook
}

func NewWebhook() Webhook {
	return &webhook{personaWebhook{}}
}

func (w webhook) Handle(ctx context.Context, data []byte, webhook WebhookType) error {
	if webhook == WebhookTypePersona {
		return w.person.Handle(ctx, data)
	}
	event := checkout.WebhookEvent{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return common.StringError(errors.Newf("error unmarshalling webhook event: %v", err))
	}
	return w.processEvent(ctx, event)
}

func (w webhook) processEvent(ctx context.Context, event checkout.WebhookEvent) error {
	switch event.Type {

	case checkout.AuthorizationApprovedEvent:
		payload := event.Data.(checkout.AuthorizationApproved)
		return w.authorizationApproved(ctx, payload)

	case checkout.AuthorizationDeclinedEvent:
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
	checkout.PostToSlack(checkout.AuthorizationApprovedEvent)
	return nil
}

func (w webhook) authorizationDeclined(ctx context.Context, data checkout.AuthorizationDeclined) error {
	log.Info().Msgf("authorization declined: %v", data)
	checkout.PostToSlack(checkout.AuthorizationDeclinedEvent)
	return nil
}

func (w webhook) paymentApproved(ctx context.Context, data checkout.PaymentApproved) error {
	log.Info().Msgf("payment approved: %v", data)
	checkout.PostToSlack(checkout.PaymentApprovedEvent)
	return nil
}

func (w webhook) paymentCaptured(ctx context.Context, data checkout.PaymentCaptured) error {
	log.Info().Msgf("payment captured: %v", data)
	checkout.PostToSlack(checkout.PaymentCapturedEvent)
	return nil
}
