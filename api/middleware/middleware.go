package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	echoDatadog "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

func CORS() echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	})
}

func Recover() echo.MiddlewareFunc {
	return echoMiddleware.Recover()
}

func Logger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			return next(c)
		}
	}
}

func LogRequest() echo.MiddlewareFunc {
	return echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger := c.Get("logger").(*zerolog.Logger)
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
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			var claims = &service.JWTClaims{}
			t, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			c.Set("userId", claims.ID)
			return t, err
		},
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

func Tracer() echo.MiddlewareFunc {
	return echoDatadog.Middleware(echoDatadog.WithServiceName("string-api"))
}
