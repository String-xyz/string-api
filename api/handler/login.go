package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type LoginType string

const (
	loginType      = LoginType("loginType")
	loginTypeNonce = LoginType("nonceSign")
)

type Login interface {
	// NoncePayload send the user a nonce payload to be signed for authentication/login purpose
	// User must provide a valid wallet address
	NoncePayload(c echo.Context) error

	//VerifySignature receives the signed noncePaylod and verifies the signature to authenticate the user.
	VerifySignature(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type login struct {
	Service service.Auth
	Group   *echo.Group
}

func NewLogin(route *echo.Echo, service service.Auth) Login {
	return &login{service, nil}
}

func (l login) NoncePayload(c echo.Context) error {
	walletAddress := c.QueryParam("walletAddress")
	if walletAddress == "" {
		return BadRequestError(c, "walletAddress must be provided")
	}
	payload, err := l.Service.PayloadToSign(walletAddress)
	if err != nil {
		LogStringError(c, err, "login: request wallet login")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "unable process payload to sign"})
	}
	return c.JSON(http.StatusOK, payload)
}

func (l login) VerifySignature(c echo.Context) error {
	var body model.WalletSignaturePayload
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: receive wallet login bind")
		return BadRequestError(c)
	}

	jwt, err := l.Service.VerifySignedPayload(body)
	if err != nil {
		LogStringError(c, err, "login: receive wallet login")
		return InternalError(c)
	}

	// set jwt in cookie
	err = SetJWTCookie(c, jwt)
	if err != nil {
		LogStringError(c, err, "login: receive email set jwt cookie")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, jwt)
}

func (l login) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	l.Group = g
	g.Use(ms...)
	g.GET("", l.NoncePayload)
	g.POST("/sign", l.VerifySignature)
}
