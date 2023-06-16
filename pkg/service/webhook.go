package service

import "golang.org/x/net/context"

type Webhook interface {
	Handle(ctx context.Context, data any) error
}

type webhook struct{}

func NewWebhook() Webhook {
	return &webhook{}
}

func (w *webhook) Handle(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) handleAuthorizationApproved(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) handleAuthorizationDeclined(ctx context.Context, data any) error {
	return nil
}

func (w *webhook) handlePaymentApproved(ctx context.Context, data any) error {
	return nil
}
