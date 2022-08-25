package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/service"
	"github.com/labstack/echo/v4"
)

type TransactionHandler interface {
	Transact(c echo.Context) error
	Quote(c echo.Context) error
	RegisterRoutes() error
}

type transactionHandler struct {
	Router  *echo.Echo
	Service service.Transaction
}

func NewTransactionHandler(route *echo.Echo, service service.Transaction) TransactionHandler {
	return &transactionHandler{route, service}
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

func (t transactionHandler) RegisterRoutes() error {
	t.Router.POST("/transact", t.Transact)
	t.Router.POST("/transact/quote", t.Quote)
	return nil
}
