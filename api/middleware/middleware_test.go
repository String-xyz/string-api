package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	env "github.com/String-xyz/go-lib/v2/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/String-xyz/string-api/config"
)

func init() {
	env.LoadEnv(&config.Var, "../../.env")
}

func TestVerifyWebhookPayload(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		secretKey      string
		signatureKey   string
		signatureValue string
		wantErr        bool
		errMessage     string
	}{
		{
			name:         "valid signature",
			body:         `{"message": "test"}`,
			secretKey:    "test-key",
			signatureKey: "Cko-Signature",
			wantErr:      false,
		},
		{
			name:         "invalid signature",
			body:         `{"message": "test"}`,
			secretKey:    "wrong-key",
			signatureKey: "Cko-Signature",
			wantErr:      true,
			errMessage:   "Failed to verify payload",
		},
		// Add more test cases as needed.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			rec := httptest.NewRecorder()

			// Calculate signature and add it to the headers
			mac := hmac.New(sha256.New, []byte(tt.secretKey))
			mac.Write([]byte(tt.body))
			expectedMAC := mac.Sum(nil)
			req.Header.Set(tt.signatureKey, hex.EncodeToString(expectedMAC))

			c := e.NewContext(req, rec)

			middleware := VerifyWebhookPayload()

			if tt.wantErr {
				err := middleware(func(c echo.Context) error {
					return nil
				})(c)
				assert.EqualError(t, err, tt.errMessage)
			} else {
				err := middleware(func(c echo.Context) error {
					return nil
				})(c)
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to compute the MAC of a given body and secret
func ComputeMAC(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
