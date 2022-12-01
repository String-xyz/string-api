package handler

import (
	"net/http"

	httpError "github.com/String-xyz/string-api/api/common/httpError"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Login interface {
	Create(c echo.Context) error
	ReceiveEmailAuthentication(c echo.Context) error
	RequestWalletLogin(c echo.Context) error
	ReceiveWalletLogin(c echo.Context) error
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
	var body model.WalletSignaturePayload
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: receive wallet login bind")
		return httpError.InvalidPayloadError(c)
	}
	jwt, err := l.Service.Create(body)
	if err != nil {
		LogStringError(c, err, "login: receive wallet login")
		return httpError.InvalidPayloadError(c) // TODO: This error is redundant. Refactor after adding body validation
	}

	// set jwt in cookie
	err = SetJWTCookie(c, jwt)
	if err != nil {
		LogStringError(c, err, "login: create set jwt cookie")
		return httpError.InternalError(c)
	}

	return c.JSON(http.StatusOK, jwt)
}

func (l login) ReceiveEmailAuthentication(c echo.Context) error {
	token := c.QueryParam("token")
	err := l.Service.ReceiveEmailAuthentication(token)
	if err != nil {
		LogStringError(c, err, "login: receive email authentication")
		return c.JSON(http.StatusBadRequest, httpError.JSONError{Message: "Invalid Token"})
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Email Successfully Authenticated"})
}

func (l login) RequestWalletLogin(c echo.Context) error {
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: request wallet login bind")
		return httpError.BadRequestError(c)
	}
	payload, err := l.Service.RequestWalletLogin(body)
	if err != nil {
		LogStringError(c, err, "login: request wallet login")
		return httpError.InternalError(c)
	}
	return c.JSON(http.StatusOK, payload)
}

func (l login) ReceiveWalletLogin(c echo.Context) error {
	var body model.WalletSignaturePayload
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: receive wallet login bind")
		return httpError.BadRequestError(c)
	}

	jwt, err := l.Service.ReceiveWalletLogin(body)
	if err != nil {
		LogStringError(c, err, "login: receive wallet login")
		return httpError.InternalError(c)
	}

	// set jwt in cookie
	err = SetJWTCookie(c, jwt)
	if err != nil {
		LogStringError(c, err, "login: receive email set jwt cookie")
		return httpError.InternalError(c)
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
	g.POST("/request", l.RequestWalletLogin)
	g.POST("", l.ReceiveWalletLogin)
}
