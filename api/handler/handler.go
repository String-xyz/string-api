package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
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
	res, err := t.Service.Execute()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (t transaction) Quote(c echo.Context) error {
	res, err := t.Service.Quote()
	if err != nil {
		return err
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
