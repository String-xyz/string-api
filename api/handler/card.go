package handler

import (
	"net/http"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	service "github.com/String-xyz/string-api/pkg/service"
)

type Card interface {
	GetAll(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type card struct {
	Service service.Card
	Group   *echo.Group
}

func NewCard(route *echo.Echo, service service.Card) Card {
	return &card{service, nil}
}

// @Summary Get all saved cards
// @Description Get all saved cards
// @Tags Cards
// @Accept json
// @Produce json
// @Security JWT
// @Success 200 {object} []checkout.CardInstrument
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /cards [get]
func (card card) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	userId, ok := c.Get("userId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid userId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	res, err := card.Service.ListByUserId(ctx, userId, platformId)
	if err != nil && errors.Cause(err).Error() != "404 Not Found" { // Not a string error
		libcommon.LogStringError(c, err, "cards: get All")
		return httperror.Internal500(c, "Cards Service Failed")
	}

	// 200
	return c.JSON(http.StatusOK, res)
}

func (card card) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	card.Group = g
	g.Use(ms...)
	g.GET("", card.GetAll)
}
