package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestVerifyWebhookPayload(t *testing.T) {
	checkoutSecretKey := "checkout_secret_key"
	personaSecretKey := "persona_secret_key"

	// Define your test cases
	tests := []struct {
		name          string
		path          string
		signatureKey  string
		signatureName string
		secretKey     string
		expectCode    int
	}{
		{
			name:          "Test unauthorized access due to invalid signature for Checkout",
			path:          "webhooks/checkout",
			signatureKey:  "invalid_signature",
			signatureName: "Cko-Signature",
			secretKey:     checkoutSecretKey,
			expectCode:    http.StatusUnauthorized,
		},
		{
			name:          "Test successful access for Checkout",
			path:          "webhooks/checkout",
			signatureKey:  computeHmacSha256("hello", checkoutSecretKey),
			signatureName: "Cko-Signature",
			secretKey:     checkoutSecretKey,
			expectCode:    http.StatusOK,
		},
		{
			name:          "Test unauthorized access due to invalid signature for Persona",
			path:          "webhooks/persona",
			signatureKey:  "t=1629478952,v1=invalid_signature",
			signatureName: "Persona-Signature",
			secretKey:     personaSecretKey,
			expectCode:    http.StatusUnauthorized,
		},
		{
			name:          "Test successful access for Persona",
			path:          "webhooks/persona",
			signatureKey:  fmt.Sprintf("t=1629478952,v1=%s", computeHmacSha256("1629478952.hello", personaSecretKey)),
			signatureName: "Persona-Signature",
			secretKey:     personaSecretKey,
			expectCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Echo context with request and response
			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/", strings.NewReader("hello"))
			req.Header.Set(tt.signatureName, tt.signatureKey)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			// Assert middleware function
			middleware := VerifyWebhookPayload(personaSecretKey, checkoutSecretKey)
			middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "Test")
			})(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

// Utility function to compute HMAC for testing
func computeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
