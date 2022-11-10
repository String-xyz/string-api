package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Auth interface {
	NonceChallenge(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type auth struct {
	service service.Auth
	logger  *zerolog.Logger
}

func NewAuth(service service.Auth) Auth {
	return &auth{service: service}
}

func (o auth) NonceChallenge(c echo.Context) error {

	param := struct {
		PublicAddress string `param:"address"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	nonce, err := o.service.Challenge(param.PublicAddress)
	if err != nil {
		o.logger.Err(err).Msg("auth challenge")
		return echo.NewHTTPError(http.StatusInternalServerError, "Nonce Challenge Service Failed")
	}

	return c.JSON(http.StatusOK, map[string]string{"nonce": nonce})
}

func (o auth) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the auth handler")
	}
	g.Use(ms...)
	g.GET("/:address/nonce", o.NonceChallenge)
}
