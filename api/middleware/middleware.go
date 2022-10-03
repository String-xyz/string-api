package middleware

import (
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func allowOrigin(origin string) (bool, error) {
	// TODO: Modify to be more restrictive
	return regexp.MatchString(`*`, origin)
}

func CORS() echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOriginFunc: allowOrigin,
		AllowMethods:    []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	})
}

func Recover() echo.MiddlewareFunc {
	return echoMiddleware.Recover()
}

func Logger(logger zerolog.Logger) echo.MiddlewareFunc {
	return echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("requestId", v.RequestID).
				Str("host", v.Host).
				Dur("latency", time.Duration(v.Latency.Milliseconds())).
				Msg("request")
			return nil
		},
	})
}

// RequestID generates a unique request ID
func RequestID() echo.MiddlewareFunc {
	return echoMiddleware.RequestID()
}

func BearerAuth() echo.MiddlewareFunc {
	config := echoMiddleware.JWTConfig{
		Claims:     &service.JWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	}
	return echoMiddleware.JWTWithConfig(config)
}

func APIKeyAuth(service service.Auth) echo.MiddlewareFunc {
	config := echoMiddleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(auth string, c echo.Context) (bool, error) {
			valid := service.ValidateAPIKey(auth)
			return valid, nil
		},
	}
	return echoMiddleware.KeyAuthWithConfig(config)
}
