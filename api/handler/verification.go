package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Verification interface {
	Authenticate(c echo.Context) error // Takes e-mail and wallet addr of user, validates email with twilio
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type verification struct {
	Service service.User
	Group   *echo.Group
}

func NewVerification(route *echo.Echo, service service.User) Verification {
	return &verification{service, nil}
}

func (v verification) Authenticate(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	// Token was provided
	token := c.QueryParam("token")
	err := v.Service.ReceiveEmailAuthentication(token)
	if err != nil {
		lg.Err(err).Msg("user authenticate")
		return c.String(http.StatusBadRequest, "Invalid Token")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Email Validated Successfully"})
}

func (v verification) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	v.Group = g
	g.Use(ms...)
	g.POST("/email", v.Authenticate)
}
