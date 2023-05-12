package handler

import (
	"net/http"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Verification interface {
	Verify(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type verification struct {
	service       service.Verification
	deviceService service.Device
	group         *echo.Group
}

func NewVerification(route *echo.Echo, service service.Verification, deviceService service.Device) Verification {
	return &verification{service, deviceService, nil}
}

// @Summary Verify
// @Description Verify receives payload from an email link sent previsouly.
// @Tags Verification
// @Accept json
// @Produce json
// @Param type query string true "type"
// @Param token query string true "token"
// @Success 200 {object} ResultMessage
// @Failure 400 {object} error
// @Router /verification [get]
func (v verification) Verify(c echo.Context) error {
	verificationType := c.QueryParam("type")
	if verificationType == "" {
		return httperror.BadRequest400(c)
	}
	if verificationType == "email" {
		// ?
		return v.verifyEmail(c)
	}

	// ?
	return v.verifyDevice(c)
}

// Verify Email receives payload from an email link sent previsouly.
// It creates a contact and sets the email as verified.
// This is a public endpoint since is called outside of a platform
func (v verification) verifyEmail(c echo.Context) error {
	ctx := c.Request().Context()

	token := c.QueryParam("token")
	err := v.service.VerifyEmailWithEncryptedToken(ctx, token)
	if err != nil {
		libcommon.LogStringError(c, err, "verification: email verification")
		return httperror.BadRequest400(c)
	}
	// 200
	return c.JSON(http.StatusOK, ResultMessage{Status: "Email successfully verified"})
}

// Verify Device receives payload from an email link sent in the login/sign.
func (v verification) verifyDevice(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.QueryParam("token")
	err := v.deviceService.VerifyDevice(ctx, token)
	if err != nil {
		libcommon.LogStringError(c, err, "verification: device verification")
		return httperror.BadRequest400(c)
	}
	// 200
	return c.JSON(http.StatusOK, ResultMessage{Status: "Device successfully verified"})
}

func (v verification) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the verification handler")
	}
	v.group = g
	g.Use(ms...)
	g.GET("", v.Verify)
}
