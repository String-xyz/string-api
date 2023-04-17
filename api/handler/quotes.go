package handler

import (
	"net/http"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// TODO: add these to stringerror in go-lib
var FUNC_NOT_ALLOWED = errors.New("function is not allowed on this contract")
var CONTRACT_NOT_ALLOWED = errors.New("contract not allowed by platform on network")

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

func (q quote) Quote(c echo.Context) error {
	ctx := c.Request().Context()
	var body model.TransactionRequest
	err := c.Bind(&body) // 'tag' binding: struct fields are annotated
	if err != nil {
		libcommon.LogStringError(c, err, "quote: quote bind")
		return httperror.BadRequestError(c)
	}
	SanitizeChecksums(&body.CxAddr, &body.UserAddress)
	// Sanitize Checksum for body.CxParams?  It might look like this:
	for i := range body.CxParams {
		SanitizeChecksums(&body.CxParams[i])
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "Platform ID not found")
	}

	res, err := q.Service.Quote(ctx, body, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "quote: quote")

		if errors.Cause(err).Error() == "w3: response handling failed: execution reverted" { // TODO: use a custom error
			return httperror.BadRequestError(c, "The requested blockchain operation will revert")
		}

		if serror.Is(err, FUNC_NOT_ALLOWED, CONTRACT_NOT_ALLOWED) {
			return httperror.ForbiddenError(c, "The requested blockchain operation is not allowed")
		}

		return httperror.InternalError(c, "Quote Service Failed")
	}

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
