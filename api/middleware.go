package api

import (
	"net/http"
	"regexp"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func allowOrigin(origin string) (bool, error) {
	// TODO: Modify to be more restrictive
	return regexp.MatchString(`*`, origin)
}

type MiddlewareConfig struct {
	any any
}

func StringMiddleware(config MiddlewareConfig) {
	echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOriginFunc: allowOrigin,
		AllowMethods:    []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	})
}
