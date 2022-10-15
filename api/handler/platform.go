package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Platform interface {
	Create(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type platform struct {
	service service.Platform
}

func NewPlatform(service service.Platform) Platform {
	return &platform{service: service}
}

func (p platform) Create(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	body := service.CreatePlatform{}
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	m, err := p.service.Create(body)
	if err != nil {
		lg.Err(err).Msg("platform create")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, m)
}

func (p platform) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the platform handler")
	}
	g.Use(ms...)
	g.POST("", p.Create)
}
