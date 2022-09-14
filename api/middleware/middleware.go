package middleware

import (
	"net/http"
	"os"
	"regexp"

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

func Logger() echo.MiddlewareFunc {
	logger := zerolog.New(os.Stdout)
	return echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("Status", v.Status).
				Str("RequestID", v.RequestID).
				Dur("latency", v.Latency).
				Msg("request")
			return nil
		},
	})
}
