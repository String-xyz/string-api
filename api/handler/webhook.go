package handler

import (
	"io"

	"github.com/labstack/echo/v4"

	"github.com/String-xyz/go-lib/v2/httperror"

	"github.com/String-xyz/string-api/pkg/service"
)

type Webhook interface {
	HandleCheckout(c echo.Context) error
	HandlePersona(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type webhook struct {
	service service.Webhook
	group   *echo.Group
}

func NewWebhook(route *echo.Echo, service service.Webhook) Webhook {
	return &webhook{service, nil}
}

func (w webhook) HandleCheckout(c echo.Context) error {
	cxt := c.Request().Context()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return httperror.BadRequest400(c, "Failed to read body")
	}

	err = w.service.Handle(cxt, body, service.WebhookTypeCheckout)
	if err != nil {
		return httperror.Internal500(c, "Failed to handle checkout webhook")
	}

	return nil
}

func (w webhook) HandlePersona(c echo.Context) error {
	cxt := c.Request().Context()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return httperror.BadRequest400(c, "Failed to read body")
	}

	err = w.service.Handle(cxt, body, service.WebhookTypePersona)
	if err != nil {
		return httperror.Internal500(c, "Failed to handle persona webhook")
	}

	return nil
}

func (w webhook) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	w.group = g
	g.POST("/checkout", w.HandleCheckout, ms...)
	g.POST("/persona", w.HandlePersona, ms...)
}
