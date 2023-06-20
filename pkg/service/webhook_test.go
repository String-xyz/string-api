package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/String-xyz/string-api/pkg/test/data"
	"github.com/stretchr/testify/assert"
)

func TestPaymentAuthorized(t *testing.T) {
	json := data.AuhorizationApprovedJSON
	err := NewWebhook().Handle(context.Background(), []byte(json))
	assert.NoError(t, err)
}

func TestPaymentDeclined(t *testing.T) {
	json := data.AuhorizationDeclinedJSON
	err := NewWebhook().Handle(context.Background(), []byte(json))
	assert.NoError(t, err)
}

func TestPaymentCaptured(t *testing.T) {
	json := data.PaymentCapturedJSON
	err := NewWebhook().Handle(context.Background(), []byte(json))
	assert.NoError(t, err)
}

func TestPaymentApproved(t *testing.T) {
	json := data.PaymentApprovedJSON
	err := NewWebhook().Handle(context.Background(), []byte(json))
	assert.NoError(t, err)
}

// Helper function to compute the MAC of a given payload and secret
func computeMAC(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
