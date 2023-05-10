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

// @Summary Transact
// @Description Transact executes a transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.ExecutionRequest true "Execution Request"
// @Success 200 {object} model.TransactionReceipt
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 500 {object} error
// @Router /transaction [post]
func (t transaction) Transact(c echo.Context) error {
	ctx := c.Request().Context()
	userId, ok := c.Get("userId").(string)
	if !ok {
		// 500
		return httperror.InternalError(c, "missing or invalid userId")
	}

	deviceId, ok := c.Get("deviceId").(string)
	if !ok {
		// 500
		return httperror.InternalError(c, "missing or invalid deviceId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		// 500
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	var body model.ExecutionRequest

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute bind")
		// 400
		return httperror.BadRequestError(c)
	}

	err = c.Validate(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute validate")
		// 400
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
			// 403
			return httperror.Unprocessable(c)
		}

		// 500
		return httperror.InternalError(c)
	}

	// 200
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
