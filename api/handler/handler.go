package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Transaction interface {
	Transact(c echo.Context) error
	Quote(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type transaction struct {
	Service service.Transaction
	Group   *echo.Group
}

func NewTransaction(route *echo.Echo, service service.Transaction) Transaction {
	return &transaction{service, nil}
}

func (t transaction) Transact(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.ExecutionRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	userId := c.Get("userId").(string)
	res, err := t.Service.Execute(body, userId) // TODO: pass in userId and use it
	if err != nil {
		lg.Err(err).Msg("transaction execute")
		return c.String(http.StatusOK, "Execute Service Failed")
	}
	return c.JSON(http.StatusOK, res)
}

func (t transaction) Quote(c echo.Context) error {
	var body model.TransactionRequest
	lg := c.Get("logger").(*zerolog.Logger)
	err := c.Bind(&body) // 'tag' binding: struct fields are annotated
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	// userId := c.Get("userId").(string)
	res, err := t.Service.Quote(body) // TODO: pass in userId and use it
	if err != nil {
		lg.Err(err).Msg("transaction quote")
		return c.String(http.StatusOK, "Quote Service Failed")
	}
	return c.JSON(http.StatusOK, res)
}

func (t transaction) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the Transaction Handler")
	}
	t.Group = g
	g.Use(ms...)
	g.POST("/", t.Transact)
	g.POST("/quote", t.Quote)
}
