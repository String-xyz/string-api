package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthAPIKey interface {
	Create(c echo.Context) error
	Approve(c echo.Context) error
	List(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type authAPIKey struct {
	isInternal bool
	service    service.APIKeyStrategy
	logger     *zerolog.Logger
}

func NewAuthAPIKey(service service.APIKeyStrategy, internal bool) AuthAPIKey {
	return &authAPIKey{service: service, isInternal: internal}
}

func (o authAPIKey) Create(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	key, err := o.service.Create()
	if err != nil {
		lg.Err(err).Stack().Msg("authKey approve:create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to process request")
	}
	return c.JSON(http.StatusOK, map[string]string{"apiKey": key})
}

func (o authAPIKey) List(c echo.Context) error {
	if !o.isInternal {
		return c.String(http.StatusMethodNotAllowed, "Not Allowed")
	}
	lg := c.Get("logger").(*zerolog.Logger)
	body := struct {
		Status string `query:"status"`
	}{}
	err := c.Bind(&body)
	if err != nil {
		lg.Err(err).Stack().Msg("authKeys list: bind")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	list, err := o.service.List(100, 0, body.Status)
	if err != nil {
		lg.Err(err).Stack().Msg("authKeys list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Register Service Failed")
	}
	return c.JSON(http.StatusCreated, list)
}

func (o authAPIKey) Approve(c echo.Context) error {
	if !o.isInternal {
		return c.String(http.StatusMethodNotAllowed, "Not Allowed")
	}
	lg := c.Get("logger").(*zerolog.Logger)
	params := struct {
		ID string `query:"id"`
	}{}
	err := c.Bind(&params)
	if err != nil {
		lg.Err(err).Stack().Msg("authKey approve:bind")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to process request")
	}
	err = o.service.Approve(params.ID)
	if err != nil {
		lg.Err(err).Stack().Msg("authKey approve:approve")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to process request")
	}
	return c.String(http.StatusOK, "Success")
}

func (o authAPIKey) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the authKey handler")
	}
	g.Use(ms...)
	g.POST("", o.Create)
	g.GET("", o.List)
	g.POST("/:id/approve", o.Approve)
}
