package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
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
		return BadRequestError(c, "WalletAddress must be provided")
	}

	payload, err := l.Service.PayloadToSign(walletAddress)
	if err != nil {
		LogStringError(c, err, "login: request wallet login")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, map[string]string{"encodedPayload": payload})
}

func (l login) VerifySignature(c echo.Context) error {
	var body model.WalletSignaturePayloadSigned
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: binding body")
		return BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return InvalidPayloadError(c, err)
	}

	resp, err := l.Service.VerifySignedPayload(body)
	if err != nil {
		LogStringError(c, err, "login: verify signature")
		return BadRequestError(c, "Invalid Payload")
	}
	// set jwt in cookie
	err = SetJWTCookie(c, resp.JWT)
	if err != nil {
		LogStringError(c, err, "login: receive email set jwt cookie")
		return InternalError(c)
	}
	return c.JSON(http.StatusOK, resp)
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
