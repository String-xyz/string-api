package handler

import (
	b64 "encoding/base64"
	"net/http"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Login interface {
	RequestToSign(c echo.Context) error
	Login(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
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

// @Summary Request To Sign
// @Description RequestToSign sends the user a nonce payload to be signed for authentication/login purposes. User must provide a valid wallet address
// @Tags Login
// @Accept json
// @Produce json
// @Param walletAddress query string true "wallet address"
// @Success 200 {object} model.SignatureRequest
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /login [get]
func (l login) RequestToSign(c echo.Context) error {
	walletAddress := c.QueryParam("walletAddress")

	if !ethcommon.IsHexAddress(walletAddress) {
		return httperror.BadRequest400(c, "Invalid wallet address")
	}

	SanitizeChecksums(&walletAddress)

	// get nonce payload
	signatureRequest, err := l.Service.PayloadToSign(c.Request().Context(), walletAddress)
	if err != nil {
		// 500
		return DefaultErrorHandler(c, err, "login: RequestToSign")
	}

	// 200
	return c.JSON(http.StatusOK, signatureRequest)
}

// @Summary Login
// @Description Login receives the signed noncePayload and verifies the signature to authenticate the user.
// @Tags Login
// @Accept json
// @Produce json
// @Param bypassDevice query boolean false "bypass device"
// @Param payload body model.WalletSignaturePayloadSigned true "Wallet Signature Payload"
// @Success 200 {object} model.UserLoginResponse
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 422 {object} error
// @Failure 500 {object} error
// @Router /login/sign [post]
func (l login) Login(c echo.Context) error {
	ctx := c.Request().Context()
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	strBypassDevice := c.QueryParam("bypassDevice")
	bypassDevice := strBypassDevice == "true" // convert to bool. default is false

	var body model.WalletSignaturePayloadSigned
	if err := c.Bind(&body); err != nil {
		libcommon.LogStringError(c, err, "login: binding body")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayload400(c, err)
	}

	// base64 decode nonce
	decodedNonce, err := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "login: verify signature decode nonce")
		return httperror.BadRequest400(c)
	}
	body.Nonce = string(decodedNonce)

	resp, err := l.Service.VerifySignedPayload(ctx, body, platformId, bypassDevice)
	if err != nil {
		libcommon.LogStringError(c, err, "login: verify signature")

		if serror.Is(err, serror.UNKNOWN_DEVICE) {
			return httperror.Unprocessable422(c)
		}

		if serror.Is(err, serror.INVALID_DATA) {
			return httperror.BadRequest400(c, "Invalid Email")
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.BadRequest400(c, "Expired, request a new payload")
		}

		return httperror.BadRequest400(c, "Invalid Payload")
	}

	// Upsert IP address in user's device
	var claims = &model.JWTClaims{}
	_, _ = jwt.ParseWithClaims(resp.JWT.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Var.JWT_SECRET_KEY), nil
	})
	ip := c.RealIP()
	l.Device.UpsertDeviceIP(ctx, claims.DeviceId, ip)

	// set auth cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		libcommon.LogStringError(c, err, "login: unable to set auth cookies")
		return httperror.Internal500(c)
	}

	// 200
	return c.JSON(http.StatusOK, resp)
}

// @Summary Refresh Token
// @Description Refresh Token
// @Tags Login
// @Accept json
// @Produce json
// @Param walletAddress body string true "wallet address"
// @Success 200 {object} model.RefreshTokenResponse
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /login/refresh [post]
func (l login) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	var body model.RefreshTokenPayload
	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "login: binding body")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayload400(c, err)
	}

	SanitizeChecksums(&body.WalletAddress)

	cookie, err := c.Cookie("StringRefreshToken")
	if err != nil {
		libcommon.LogStringError(c, err, "RefreshToken: unable to get StringRefreshToken cookie")
		return httperror.Unauthorized401(c)
	}

	resp, err := l.Service.RefreshToken(ctx, cookie.Value, body.WalletAddress, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "login: refresh token")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.BadRequest400(c, "wallet address not associated with this user")
		}

		return httperror.BadRequest400(c, "Invalid or expired token")
	}

	// set auth in cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		libcommon.LogStringError(c, err, "RefreshToken: unable to set auth cookies")
		return httperror.Internal500(c)
	}

	// 200
	return c.JSON(http.StatusOK, resp)
}

// @Summary Logout
// @Description Logout of the application, invalidating the auth cookies
// @Tags Login
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /login/logout [post]
func (l login) Logout(c echo.Context) error {
	// get refresh token from cookie
	cookie, err := c.Cookie("StringRefreshToken")
	if err != nil {
		libcommon.LogStringError(c, err, "Logout: unable to get StringRefreshToken cookie")
		return httperror.Unauthorized401(c)
	}

	// invalidate refresh token. Returns error if token is not found
	err = l.Service.InvalidateRefreshToken(cookie.Value)
	if err != nil {
		libcommon.LogStringError(c, err, "Token not found")
		// if error continue anyway, at least delete the cookies
	}

	// delete auth cookies
	err = DeleteAuthCookies(c)
	if err != nil {
		libcommon.LogStringError(c, err, "Logout: unable to delete auth cookies")
		return httperror.Internal500(c)
	}

	// 204
	return c.JSON(http.StatusNoContent, nil)
}

func (l login) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	l.Group = g
	g.GET("", l.RequestToSign, ms...)
	g.POST("/sign", l.Login, ms...)
	g.POST("/refresh", l.RefreshToken, ms...)
	g.POST("/logout", l.Logout)
}
