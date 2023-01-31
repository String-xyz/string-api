package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Verification interface {
	// VerifyEmail receives payload from an email link sent previsouly
	// it creates a contact and sets the email as verified
	// this is a public endpoint since is called outside of a platform
	VerifyEmail(c echo.Context) error

	//VerifyDevice receives a payload from the email link sent from login/sign
	VerifyDevice(c echo.Context) error
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

func (v verification) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	err := v.service.VerifyEmail(token)
	if err != nil {
		LogStringError(c, err, "verification: email verification")
		return BadRequestError(c)
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Email successfully verified"})
}

func (v verification) VerifyDevice(c echo.Context) error {
	token := c.QueryParam("token")
	err := v.deviceService.VerifyDevice(token)
	if err != nil {
		LogStringError(c, err, "verification: device verification")
		return BadRequestError(c)
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Device successfully verified"})
}

func (v verification) verify(c echo.Context) error {
	verificationType := c.QueryParam("type")
	if verificationType == "" {
		return BadRequestError(c)
	}
	if verificationType == "email" {
		return v.VerifyEmail(c)
	}

	return v.VerifyDevice(c)
}

func (v verification) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the verification handler")
	}
	v.group = g
	g.Use(ms...)
	g.GET("", v.verify)
}
