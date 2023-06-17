package service

import "golang.org/x/net/context"

type Webhook interface {
	Handle(ctx context.Context, data []byte) error
}

type webhook struct{}

func NewWebhook() Webhook {
	return &webhook{}
}

func (w *webhook) Handle(ctx context.Context, data []byte) error {
	return nil
}

func (w *webhook) authorizationApproved(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) authorizationDeclined(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) paymentApproved(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) paymentCaptured(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) validatePayload(payload any) error {
	return nil
}
