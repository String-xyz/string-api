package handler

import (
	"net/http"

	commonlib "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
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
	key, err := o.service.Create()
	if err != nil {
		commonlib.LogStringError(c, err, "authKey approve: create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to process request")
	}
	return c.JSON(http.StatusOK, key)
}

func (o authAPIKey) List(c echo.Context) error {
	if !o.isInternal {
		return httperror.NotAllowedError(c)
	}
	body := struct {
		Status string `query:"status"`
		Limit  int    `query:"limit"`
		Offset int    `query:"offset"`
	}{}
	err := c.Bind(&body)
	if err != nil {
		commonlib.LogStringError(c, err, "authKey list: bind")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	list, err := o.service.List(body.Limit, body.Offset, body.Status)
	if err != nil {
		commonlib.LogStringError(c, err, "authKey list")
		return echo.NewHTTPError(http.StatusInternalServerError, "ApiKey Service Failed")
	}
	return c.JSON(http.StatusCreated, list)
}

func (o authAPIKey) Approve(c echo.Context) error {
	if !o.isInternal {
		return httperror.NotAllowedError(c)
	}
	params := struct {
		Id string `param:"id"`
	}{}
	err := c.Bind(&params)

	if err != nil {
		commonlib.LogStringError(c, err, "authKey approve: bind")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to process request")
	}
	err = o.service.Approve(params.Id)
	if err != nil {
		commonlib.LogStringError(c, err, "authKey approve: approve")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to process request")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Success"})
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
