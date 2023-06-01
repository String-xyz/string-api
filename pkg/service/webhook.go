package service

import "context"

type Webhook interface{}
type webhook struct{}

func NewWebhook() Webhook {
	return &webhook{}
}

func (w *webhook) Handle(ctx context.Context, req interface{}) error {
	return nil
}
