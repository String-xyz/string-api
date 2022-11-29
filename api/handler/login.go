package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/api/validator"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Login interface {
	Create(c echo.Context) error
	ReceiveEmailAuthentication(c echo.Context) error
	RequestEmailLogin(c echo.Context) error
	ReceiveEmailLogin(c echo.Context) error
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
	var body model.UserRequest
	if err := c.Bind(&body); err != nil {
		LogStringError(c, err, "login: create bind")
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, validator.ExtractErrorParams(err))
	}
	jwt, err := l.Service.Create(body)
	if err != nil {
		LogStringError(c, err, "login: create")
		return c.String(http.StatusBadRequest, "User Service Failed")
	}
	return c.JSON(http.StatusOK, jwt)
}

func (l login) ReceiveEmailAuthentication(c echo.Context) error {
	token := c.QueryParam("token")
	err := l.Service.ReceiveEmailAuthentication(token)
	if err != nil {
		LogStringError(c, err, "login: receive email authentication")
		return c.String(http.StatusBadRequest, "Invalid Token")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Email Successfully Authenticated"})
}

func (l login) RequestEmailLogin(c echo.Context) error {
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: request email login bind")
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = l.Service.RequestEmailLogin(body)
	if err != nil {
		LogStringError(c, err, "login: request email login")
		return c.String(http.StatusBadRequest, "User Service Failed")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "User Login Sent to Email"})
}

func (l login) ReceiveEmailLogin(c echo.Context) error {
	// Token was provided
	token := c.QueryParam("token")
	jwt, err := l.Service.ReceiveEmailLogin(token)
	if err != nil {
		LogStringError(c, err, "login: receive email login")
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
	g.GET("/email", l.ReceiveEmailAuthentication)
	g.POST("/new", l.Create)
	g.POST("/request", l.RequestEmailLogin)
	g.GET("", l.ReceiveEmailLogin)
}
