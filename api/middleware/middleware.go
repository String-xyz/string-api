package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
)

func JWTAuth() echo.MiddlewareFunc {
	config := echoMiddleware.JWTConfig{
		TokenLookup: "header:Authorization,cookie:StringJWT",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			var claims = &model.JWTClaims{}
			t, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.Var.JWT_SECRET_KEY), nil
			})

			c.Set("userId", claims.UserId)
			c.Set("deviceId", claims.DeviceId)
			c.Set("platformId", claims.PlatformId)

			return t, err
		},
		SigningKey: []byte(config.Var.JWT_SECRET_KEY),
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			libcommon.LogStringError(c, err, "Error in JWTAuth middleware")

			return httperror.Unauthorized401(c)
		},
	}
	return echoMiddleware.JWTWithConfig(config)
}

func APIKeyPublicAuth(service service.Auth) echo.MiddlewareFunc {
	config := echoMiddleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(auth string, c echo.Context) (bool, error) {
			platformId, err := service.ValidateAPIKeyPublic(c.Request().Context(), auth)
			if err != nil {
				libcommon.LogStringError(c, err, "Error in APIKeyPublicAuth middleware")
				return false, err
			}

			c.Set("platformId", platformId)

			return true, nil
		},
	}
	return echoMiddleware.KeyAuthWithConfig(config)
}

func APIKeySecretAuth(service service.Auth) echo.MiddlewareFunc {
	config := echoMiddleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(auth string, c echo.Context) (bool, error) {
			platformId, err := service.ValidateAPIKeySecret(c.Request().Context(), auth)
			if err != nil {
				libcommon.LogStringError(c, err, "Error in APIKeySecretAuth middleware")
				return false, err
			}

			// TODO: Validate platformId
			c.Set("platformId", platformId)
			return true, nil
		},
	}
	return echoMiddleware.KeyAuthWithConfig(config)
}

func Georestrict(service service.Geofencing) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			/* Get Ip from request */
			ip := c.RealIP()

			// check if location ip is restricted
			isAllowed, err := service.IsAllowed(ip)
			// in case of error, what we should do? allow or deny?
			// For now we are denying
			if err != nil || !isAllowed {
				if err != nil {
					libcommon.LogStringError(c, err, "Error in georestrict middleware")
				}
				return httperror.Forbidden403(c, "Error: Geo Location Forbidden")
			}

			return next(c)
		}
	}
}

func VerifyWebhookPayload(pskey string, ckoskey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var signatureHeaderName string
			var secretKey string
			var validateFunc func([]byte, string, string) bool

			switch c.Path() {
			case "webhooks/checkout":
				signatureHeaderName = "Cko-Signature"
				secretKey = ckoskey
				validateFunc = validateSignatureCheckout
			case "webhooks/persona":
				signatureHeaderName = "Persona-Signature"
				secretKey = pskey
				validateFunc = validateSignaturePersona
			default:
				return httperror.BadRequest400(c, "Invalid path")
			}

			signatureHeader := c.Request().Header.Get(signatureHeaderName)
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return httperror.BadRequest400(c, "Failed to read body")
			}

			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

			if !validateFunc(body, signatureHeader, secretKey) {
				return httperror.Unauthorized401(c, "Failed to verify payload")
			}

			return next(c)
		}
	}
}

func validateSignatureCheckout(body []byte, signature string, secretKey string) bool {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)

	receivedMAC, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(receivedMAC, expectedMAC)
}

func validateSignaturePersona(body []byte, signatureHeader string, secretKey string) bool {
	parts := strings.Split(signatureHeader, ",")
	var timestamp, signature string
	for _, part := range parts {
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
		} else if strings.HasPrefix(part, "v1=") {
			signature = strings.TrimPrefix(part, "v1=")
		}
	}

	macData := timestamp + "." + string(body)

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(macData))
	expectedMAC := mac.Sum(nil)

	receivedMAC, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedMAC, receivedMAC)
}
