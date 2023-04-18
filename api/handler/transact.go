package handler

import (
	"net/http"
	"strings"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
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
	ctx := c.Request().Context()
	userId, ok := c.Get("userId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid userId")
	}

	deviceId, ok := c.Get("deviceId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid deviceId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	var body model.ExecutionRequest

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute bind")
		return httperror.BadRequestError(c)
	}

	err = c.Validate(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute validate")
		return httperror.InvalidPayloadError(c, err)
	}

	transactionRequest := body.Quote.TransactionRequest

	SanitizeChecksums(&transactionRequest.CxAddr, &transactionRequest.UserAddress)
	// Sanitize Checksum for body.CxParams?  It might look like this:
	for i := range transactionRequest.CxParams {
		SanitizeChecksums(&transactionRequest.CxParams[i])
	}

	ip := c.RealIP()

	res, err := t.Service.Execute(ctx, body, userId, deviceId, platformId, ip)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute")

		if strings.Contains(err.Error(), "risk:") || strings.Contains(err.Error(), "payment:") {
			return httperror.Unprocessable(c)
		}

		return httperror.InternalError(c)
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
