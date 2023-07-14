package handler

import (
	"net/http"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Quotes interface {
	Quote(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type quote struct {
	Service service.Transaction
	Group   *echo.Group
}

func NewQuote(route *echo.Echo, service service.Transaction) Quotes {
	return &quote{service, nil}
}

// @Summary Quote
// @Description Quote returns the estimated cost of a transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.TransactionRequest true "Transaction Request"
// @Success 200 {object} model.Quote
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 500 {object} error
// @Router /quote [post]
func (q quote) Quote(c echo.Context) error {
	ctx := c.Request().Context()
	var body model.TransactionRequest

	err := c.Bind(&body) // 'tag' binding: struct fields are annotated
	if err != nil {
		libcommon.LogStringError(c, err, "quote: quote bind")
		return httperror.BadRequest400(c)
	}

	err = c.Validate(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "quote: quote validate")
		return httperror.InvalidPayload400(c, err)
	}

	// TODO: See if there's a way to batch these into a single call of SanitizeChecksums
	SanitizeChecksums(&body.UserAddress)
	for i := range body.Actions {
		SanitizeChecksums(&body.Actions[i].CxAddr, &body.UserAddress)
		for j := range body.Actions[i].CxParams {
			SanitizeChecksums(&body.Actions[i].CxParams[j])
		}
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	res, err := q.Service.Quote(ctx, body, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "quote: quote")

		if errors.Cause(err).Error() == "w3: response handling failed: execution reverted" { // TODO: use a custom error
			return httperror.BadRequest400(c, "The requested blockchain operation will revert")
		}

		if serror.Is(err, serror.FUNC_NOT_ALLOWED, serror.CONTRACT_NOT_ALLOWED) {
			return httperror.Forbidden403(c, "The requested blockchain operation is not allowed")
		}

		return httperror.Internal500(c, "Quote Service Failed")
	}

	// 200
	return c.JSON(http.StatusOK, res)
}

func (q quote) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the Quote Handler")
	}
	q.Group = g
	g.Use(ms...)
	g.POST("", q.Quote)
}
