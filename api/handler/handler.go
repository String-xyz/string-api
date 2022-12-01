package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/api/common/httpError"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
	var body model.ExecutionRequest
	err := c.Bind(&body)
	if err != nil {
                LogStringError(c, err, "transact: execute bind")
		return httpError.BadRequestError(c)
	}
	userId := c.Get("userId").(string)
	res, err := t.Service.Execute(body, userId)
	if err != nil {
		LogStringError(c, err, "transact: execute")
		return httpError.InternalError(c)
	}
	return c.JSON(http.StatusOK, res)
}

func (t transaction) Quote(c echo.Context) error {
	var body model.TransactionRequest
	err := c.Bind(&body) // 'tag' binding: struct fields are annotated
	if err != nil {
		return httpError.BadRequestError(c)
	}
	// userId := c.Get("userId").(string)
	res, err := t.Service.Quote(body) // TODO: pass in userId and use it
	if err != nil && errors.Cause(err).Error() == "w3: response handling failed: execution reverted" {
		return c.JSON(http.StatusBadRequest, httpError.JSONError{Message: "The requested blockchain operation will revert"})
	} else if err != nil {
		LogStringError(c, err, "transact: quote")
		return c.JSON(http.StatusInternalServerError, httpError.JSONError{Message: "Quote Service Failed"})
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
