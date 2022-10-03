package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	NonceChallenge(c echo.Context) error
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, jwt)
}

func (o auth) Login(c echo.Context) error {
	body := struct {
		service.UserPKLogin
		service.UserLoginEmail
		LoginType string `query:"loginType"`
	}{}

	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	switch body.LoginType {
	case "email":
		return o.LoginEmail(c, body.UserLoginEmail)
	case "privateKey":
		return o.LoginPK(c, body.UserPKLogin)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid login type")
	}
}

func (o auth) LoginEmail(c echo.Context, body service.UserLoginEmail) error {
	jwt, err := o.service.LoginEmail(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jwt)
}

func (o auth) LoginPK(c echo.Context, body service.UserPKLogin) error {
	jwt, err := o.service.LoginPK(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jwt)
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"nonce": nonce})
}

func (o auth) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the auth handler")
	}
	g.Use(ms...)
	g.POST("/register", o.Register)
	g.POST("/login", o.Login)
	g.GET("/:address/nonce", o.NonceChallenge)
}
