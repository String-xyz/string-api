package handler

import (
	"net/http"

	libCommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
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
	body := service.CreatePlatform{}
	err := c.Bind(&body)
	if err != nil {
		libCommon.LogStringError(c, err, "platform: create bind")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	m, err := p.service.Create(body)
	if err != nil {
		libCommon.LogStringError(c, err, "platform: create")
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
