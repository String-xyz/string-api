package handler

import (
	"net/http"
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
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
// @Security JWT
// @Param saveCard query boolean false "do not save payment info"
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
		return httperror.Internal500(c, "missing or invalid userId")
	}

	deviceId, ok := c.Get("deviceId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid deviceId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	var body model.ExecutionRequest

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute bind")
		return httperror.BadRequest400(c)
	}

	err = c.Validate(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute validate")
		return httperror.InvalidPayload400(c, err)
	}

	transactionRequest := body.Quote.TransactionRequest

	// TODO: These should already be sanitized by the quote, double check when there's time
	SanitizeChecksums(&transactionRequest.UserAddress)
	for i := range transactionRequest.Actions {
		SanitizeChecksums(&transactionRequest.Actions[i].CxAddr, &transactionRequest.UserAddress)
		for j := range transactionRequest.Actions[i].CxParams {
			SanitizeChecksums(&transactionRequest.Actions[i].CxParams[j])
		}
	}

	ip := c.RealIP()

	res, err := t.Service.Execute(ctx, body, userId, deviceId, platformId, ip)
	if err != nil {
		libcommon.LogStringError(c, err, "transact: execute")

		if strings.Contains(err.Error(), "risk:") || strings.Contains(err.Error(), "payment:") {
			return httperror.Unprocessable422(c)
		}

		return httperror.Internal500(c)
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
