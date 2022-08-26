package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/service"
	"github.com/labstack/echo/v4"
)

type TransactionHandler interface {
	Transact(c echo.Context) error
	Quote(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type transactionHandler struct {
	Service service.Transaction
	Group   *echo.Group
}

func NewTransactionHandler(route *echo.Echo, service service.Transaction) TransactionHandler {
	return &transactionHandler{service, nil}
}

func (t transactionHandler) Transact(c echo.Context) error {
	res, err := t.Service.Execute()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (t transactionHandler) Quote(c echo.Context) error {
	res, err := t.Service.Quote()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (t transactionHandler) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to TransactionHandler")
	}
	t.Group = g
	g.Use(ms...)
	g.POST("/", t.Transact)
	g.POST("/quote", t.Quote)
}
