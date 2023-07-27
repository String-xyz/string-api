package service

import (
	"context"
	"encoding/json"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"

	"github.com/String-xyz/string-api/pkg/internal/persona"
)

type personaWebhook struct{}

func (p personaWebhook) Handle(ctx context.Context, data []byte) error {
	event := persona.Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return common.StringError(errors.Newf("error unmarshalling webhook event: %v", err))
	}
	return p.processEvent(ctx, event)
}

func (p personaWebhook) processEvent(ctx context.Context, event persona.Event) error {
	payload, err := event.Attributes.GetPayloadData()
	if err != nil {
		return common.StringError(errors.Newf("error getting payload data: %v", err))
	}
	switch event.Attributes.Name {
	case persona.EventTypeAccountCreated:
		return p.account(ctx, payload, event.Attributes.Name)
	case persona.EventTypeInquiryCreated, persona.EventTypeInquiryStarted, persona.EventTypeInquiryCompleted:
		return p.inquiry(ctx, payload, event.Attributes.Name)
	case persona.EventTypeVerificationCreated, persona.EventTypeVerificationPassed, persona.EventTypeVerificationFailed:
		return p.verification(ctx, payload, event.Attributes.Name)
	default:
		return common.StringError(errors.Newf("unknown event type: %s", event.Attributes.Name))
	}
}

func (p personaWebhook) account(ctx context.Context, payload persona.PayloadData, eventType persona.EventType) error {
	account, ok := payload.(persona.Account)
	if !ok {
		return common.StringError(errors.New("error casting payload to account"))
	}
	log.Info().Interface("account", account).Msg("account event")
	return nil
}

func (p personaWebhook) inquiry(ctx context.Context, payload persona.PayloadData, eventType persona.EventType) error {
	inquiry, ok := payload.(persona.Inquiry)
	if !ok {
		return common.StringError(errors.New("error casting payload to inquiry"))
	}
	log.Info().Interface("inquiry", inquiry).Msg("inquiry event")
	return nil
}

func (p personaWebhook) verification(ctx context.Context, payload persona.PayloadData, eventType persona.EventType) error {
	verification, ok := payload.(persona.Verification)
	if !ok {
		return common.StringError(errors.New("error casting payload to verification"))
	}
	log.Info().Interface("verification", verification).Msg("verification event")
	return nil
}
