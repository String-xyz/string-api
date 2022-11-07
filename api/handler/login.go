package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Login interface {
	Authenticate(c echo.Context) error // Takes e-mail and wallet addr of user, validates email with twilio
	Create(c echo.Context) error
	Request(c echo.Context) error
	AuthenticateLogin(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type login struct {
	Service service.User
	Group   *echo.Group
}

func NewLogin(route *echo.Echo, service service.User) Login {
	return &login{service, nil}
}

func (l login) Create(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = l.Service.Create(body)
	if err != nil {
		lg.Err(err).Msg("user create")
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "User Authentication Sent to Email"})
}

func (l login) Authenticate(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	// Token was provided
	token := c.QueryParam("token")
	jwt, err := l.Service.ReceiveEmailAuthentication(token)
	if err != nil {
		lg.Err(err).Msg("user authenticate")
		return c.String(http.StatusBadRequest, "Invalid Token")
	}
	return c.JSON(http.StatusOK, jwt)
}

func (l login) Request(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = l.Service.RequestEmailLogin(body)
	if err != nil {
		lg.Err(err).Msg("user login")
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "User Login Sent to Email"})
}

func (l login) AuthenticateLogin(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	// Token was provided
	token := c.QueryParam("token")
	jwt, err := l.Service.ReceiveEmailLogin(token)
	if err != nil {
		lg.Err(err).Msg("user authenticate")
		return c.String(http.StatusBadRequest, "Invalid Token")
	}
	return c.JSON(http.StatusOK, jwt)
}

func (l login) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	l.Group = g
	g.Use(ms...)
	g.POST("/email", l.Authenticate)
	g.POST("/new", l.Create)
	g.POST("/request", l.Request)
	g.POST("", l.AuthenticateLogin)
}
