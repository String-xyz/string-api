package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/String-xyz/string-api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.LoadEnv()
}

func TestVerifyWebhookPayload(t *testing.T) {
	secretKey := config.Var.WEBHOOK_SECRET_KEY

	// We'll test with two cases
	tests := []struct {
		name         string
		giveBody     []byte
		giveMAC      string
		wantHTTPCode int
	}{
		{
			// This test case provides a valid body and MAC
			name:         "Valid MAC",
			giveBody:     []byte("Hello, World!"),
			giveMAC:      ComputeMAC([]byte("Hello, World!"), secretKey),
			wantHTTPCode: http.StatusOK,
		},
		{
			// This test case provides an invalid MAC
			name:         "Invalid MAC",
			giveBody:     []byte("Hello, World!"),
			giveMAC:      ComputeMAC([]byte("Bye, World!"), secretKey),
			wantHTTPCode: http.StatusUnauthorized,
		},
	}

	// Let's iterate over our test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Instantiate Echo
			e := echo.New()

			// Our middleware under test
			middleware := VerifyWebhookPayload()

			// Mock a request
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.giveBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Signature", tt.giveMAC)

			// Mock a response recorder
			rec := httptest.NewRecorder()

			// Create a context for our request
			c := e.NewContext(req, rec)

			// Mock a next function
			next := func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			}

			// Call our middleware
			err := middleware(next)(c)

			// There should be no error returned
			assert.NoError(t, err)

			// Check if the status code is what we expect
			assert.Equal(t, tt.wantHTTPCode, rec.Code)
		})
	}
}

// Helper function to compute the MAC of a given body and secret
func ComputeMAC(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
