package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type auth struct {
	service service.Auth
}

func NewAuth(service service.Auth) Auth {
	return &auth{service: service}
}

func (o auth) Register(c echo.Context) error {
	var body service.UserRegister
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	jwt, err := o.service.Register(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, jwt)
}

func (o auth) Login(c echo.Context) error {
	var body service.UserLogin
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	jwt, err := o.service.LoginEmail(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, jwt)
}

func (o auth) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the Transaction Handler")
	}
	g.Use(ms...)
	g.POST("/register", o.Register)
	g.POST("/login", o.Login)
}
