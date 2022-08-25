package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/service"
	"github.com/labstack/echo/v4"
)

type TransactionHandler interface {
	Transact(c echo.Context) error
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
	return c.String(http.StatusOK, "Hello, World!")
}

func (t transactionHandler) RegisterRoutes() error {
	t.Router.GET("/", t.Transact)
	return nil
}
