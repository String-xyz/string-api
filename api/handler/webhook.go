package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/String-xyz/string-api/pkg/service"
)

type Webhook interface {
	Handle(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type webhook struct {
	Service service.Webhook
	Group   *echo.Group
}

func NewWebhook(route *echo.Echo, service service.Webhook) Webhook {
	return &webhook{service, nil}
}

func (w *webhook) Handle(c echo.Context) error {
	return nil
}

func (w *webhook) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	g.POST("/checkout", w.Handle, ms...)
}
