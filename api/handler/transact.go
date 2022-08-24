package handler

import "github.com/labstack/echo/v4"

type TransactionHandler interface {
	Transact(c echo.Context) error
	RegisterRoutes() error
}

type transactionHandler struct {
	Router *echo.Echo
}

func NewTransactionHandler(route *echo.Echo) TransactionHandler {
	return &transactionHandler{route}
}

func (t transactionHandler) Transact(c echo.Context) error {
	return nil
}

func (t transactionHandler) RegisterRoutes() error {
	t.Router.GET("/", t.Transact)
	return nil
}
