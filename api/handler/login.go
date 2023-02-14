package handler

import (
	b64 "encoding/base64"
	"net/http"
	"strings"

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
	RefreshToken(c echo.Context) error
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
	SanitizeChecksums(&walletAddress)
	payload, err := l.Service.PayloadToSign(walletAddress)
	if err != nil {
		LogStringError(c, err, "login: request wallet login")
		return InternalError(c)
	}

	encodedNonce := b64.StdEncoding.EncodeToString([]byte(payload.Nonce))
	return c.JSON(http.StatusOK, map[string]string{"nonce": encodedNonce})
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

	// base64 decode nonce
	decodedNonce, _ := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		LogStringError(c, err, "login: verify signature decode nonce")
		return BadRequestError(c)
	}
	body.Nonce = string(decodedNonce)

	resp, err := l.Service.VerifySignedPayload(body)
	if err != nil {
		if strings.Contains(err.Error(), "unknown device") {
			return Unprocessable(c)
		}
		if strings.Contains(err.Error(), "invalid email") {
			return InvalidEmail(c)
		}

		LogStringError(c, err, "login: verify signature")
		return BadRequestError(c, "Invalid Payload")
	}
	// set auth cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		LogStringError(c, err, "login: unable to set auth cookies")
		return InternalError(c)
	}
	return c.JSON(http.StatusOK, resp)
}

func (l login) RefreshToken(c echo.Context) error {
	var body model.RefreshTokenPayload
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "login: binding body")
		return BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return InvalidPayloadError(c, err)
	}

	SanitizeChecksums(&body.WalletAddress)

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		LogStringError(c, err, "RefreshToken: unable to get refresh_token cookie")
		return Unauthorized(c)
	}

	resp, err := l.Service.RefreshToken(cookie.Value, body.WalletAddress)
	if err != nil {
		if strings.Contains(err.Error(), "wallet address not associated with this user") {
			return BadRequestError(c, "wallet address not associated with this user")
		}

		LogStringError(c, err, "login: refresh token")
		return BadRequestError(c, "Invalid or expired token")
	}

	// set auth in cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		LogStringError(c, err, "RefreshToken: unable to set auth cookies")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, resp)
}

// logout
func (l login) Logout(c echo.Context) error {
	// get refresh token from cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		LogStringError(c, err, "Logout: unable to get refresh_token cookie")
		return Unauthorized(c)
	}

	// invalidate refresh token. Returns error if token is not found
	err = l.Service.InvalidateRefreshToken(cookie.Value)
	if err != nil {
		LogStringError(c, err, "Token not found")
	}
	// There is no need to invalidate the access token since it is a short lived token

	// delete auth cookies
	err = DeleteAuthCookies(c)
	if err != nil {
		LogStringError(c, err, "Logout: unable to delete auth cookies")
		return InternalError(c)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (l login) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	l.Group = g
	g.Use(ms...)
	g.GET("", l.NoncePayload)
	g.POST("/sign", l.VerifySignature)
	g.POST("/refresh", l.RefreshToken)
	g.POST("/logout", l.Logout)
}
