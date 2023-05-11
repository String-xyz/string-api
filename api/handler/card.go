package handler

import (
	"net/http"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	service "github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
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
// @Security ApiKeyAuth
// @Success 200 {object} []checkout.CustomerInstrument
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /cards [get]
func (card card) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	userId, ok := c.Get("userId").(string)
	if !ok {
		// 500
		return httperror.InternalError(c, "missing or invalid userId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		// 500
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	res, err := card.Service.FetchSavedCards(ctx, userId, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "cards: get All")
		// 500
		return httperror.InternalError(c, "Cards Service Failed")
	}

	// 200
	return c.JSON(http.StatusOK, res)
}

func (card card) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	card.Group = g
	g.Use(ms...)
	g.GET("", card.GetAll)
}
