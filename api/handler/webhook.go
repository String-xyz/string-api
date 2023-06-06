package handler

import "github.com/labstack/echo/v4"

type Webhook interface {
	Handle(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type webhook struct {
	Group *echo.Group
}

func NewWebhook(route *echo.Echo) Webhook {
	return &webhook{}
}

func (w *webhook) Handle(c echo.Context) error {
	return nil
}

func (w *webhook) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	g.POST("/checkout", w.Handle, ms...)
}
