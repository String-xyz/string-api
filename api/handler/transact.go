package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Transaction interface {
	Transact(c echo.Context) error
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
		return BadRequestError(c)
	}

	SanitizeChecksums(&body.CxAddr, &body.UserAddress)
	// Sanitize Checksum for body.CxParams?  It might look like this:
	for i := range body.CxParams {
		SanitizeChecksums(&body.CxParams[i])
	}
	userId := c.Get("userId").(string)
	deviceId := c.Get("deviceId").(string)
	res, err := t.Service.Execute(body, userId, deviceId)
	if err != nil && strings.Contains(err.Error(), "risk:") {
		LogStringError(c, err, "transact: execute")
		return Unprocessable(c)
	}
	if err != nil {
		LogStringError(c, err, "transact: execute")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, res)
}

func (t transaction) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the Transaction Handler")
	}
	t.Group = g
	g.Use(ms...)
	g.POST("", t.Transact)
}
