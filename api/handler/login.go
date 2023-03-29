package handler

import (
	b64 "encoding/base64"
	"net/http"
	"os"
	"strings"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Login interface {
	// NoncePayload send the user a nonce payload to be signed for authentication/login purpose
	// User must provide a valid wallet address
	NoncePayload(c echo.Context) error

	//VerifySignature receives the signed noncePaylod and verifies the signature to authenticate the user.
	VerifySignature(c echo.Context) error
	RegisterRoutes(g *echo.Group, apikeyMid echo.MiddlewareFunc)
	RefreshToken(c echo.Context) error
}

type login struct {
	Service service.Auth
	Device  service.Device
	Group   *echo.Group
}

func NewLogin(route *echo.Echo, service service.Auth, device service.Device) Login {
	return &login{service, device, nil}
}

func (l login) NoncePayload(c echo.Context) error {
	walletAddress := c.QueryParam("walletAddress")
	if walletAddress == "" {
		return httperror.BadRequestError(c, "WalletAddress must be provided")
	}
	SanitizeChecksums(&walletAddress)
	payload, err := l.Service.PayloadToSign(walletAddress)
	if err != nil {
		libcommon.LogStringError(c, err, "login: request wallet login")
		return httperror.InternalError(c)
	}

	encodedNonce := b64.StdEncoding.EncodeToString([]byte(payload.Nonce))
	return c.JSON(http.StatusOK, map[string]string{"nonce": encodedNonce})
}

func (l login) VerifySignature(c echo.Context) error {
	platformId := c.Get("platformId").(string)

	ctx := c.Request().Context()
	var body model.WalletSignaturePayloadSigned
	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "login: binding body")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	// base64 decode nonce
	decodedNonce, _ := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "login: verify signature decode nonce")
		return httperror.BadRequestError(c)
	}
	body.Nonce = string(decodedNonce)

	resp, err := l.Service.VerifySignedPayload(ctx, body, platformId)
	if err != nil {
		if strings.Contains(err.Error(), "unknown device") {
			return httperror.Unprocessable(c)
		}
		if strings.Contains(err.Error(), "invalid email") {
			return httperror.BadRequestError(c, "Invalid Email")
		}

		libcommon.LogStringError(c, err, "login: verify signature")
		return httperror.BadRequestError(c, "Invalid Payload")
	}

	// Upsert IP address in user's device
	var claims = &service.JWTClaims{}
	_, _ = jwt.ParseWithClaims(resp.JWT.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	ip := c.RealIP()
	l.Device.UpsertDeviceIP(ctx, claims.DeviceId, ip)

	// set auth cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		libcommon.LogStringError(c, err, "login: unable to set auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (l login) RefreshToken(c echo.Context) error {
	platformId := c.Get("platformId").(string)

	ctx := c.Request().Context()
	var body model.RefreshTokenPayload
	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "login: binding body")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	SanitizeChecksums(&body.WalletAddress)

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		libcommon.LogStringError(c, err, "RefreshToken: unable to get refresh_token cookie")
		return httperror.Unauthorized(c)
	}

	resp, err := l.Service.RefreshToken(ctx, cookie.Value, body.WalletAddress, platformId)
	if err != nil {
		if strings.Contains(err.Error(), "wallet address not associated with this user") {
			return httperror.BadRequestError(c, "wallet address not associated with this user")
		}

		libcommon.LogStringError(c, err, "login: refresh token")
		return httperror.BadRequestError(c, "Invalid or expired token")
	}

	// set auth in cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		libcommon.LogStringError(c, err, "RefreshToken: unable to set auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, resp)
}

// logout
func (l login) Logout(c echo.Context) error {
	// get refresh token from cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		libcommon.LogStringError(c, err, "Logout: unable to get refresh_token cookie")
		return httperror.Unauthorized(c)
	}

	// invalidate refresh token. Returns error if token is not found
	err = l.Service.InvalidateRefreshToken(cookie.Value)
	if err != nil {
		libcommon.LogStringError(c, err, "Token not found")
	}
	// There is no need to invalidate the access token since it is a short lived token

	// delete auth cookies
	err = DeleteAuthCookies(c)
	if err != nil {
		libcommon.LogStringError(c, err, "Logout: unable to delete auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (l login) RegisterRoutes(g *echo.Group, apikeyMid echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	l.Group = g
	g.GET("", l.NoncePayload)
	g.POST("/sign", l.VerifySignature, apikeyMid)
	g.POST("/refresh", l.RefreshToken, apikeyMid)
	g.POST("/logout", l.Logout)
}
