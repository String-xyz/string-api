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

func (card card) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get("userId").(string)
	res, err := card.Service.FetchSavedCards(ctx, userId)
	if err != nil {
		libcommon.LogStringError(c, err, "cards: cards")
		return httperror.InternalError(c, "Cards Service Failed")
	}
	return c.JSON(http.StatusOK, res)
}

func (card card) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	card.Group = g
	g.Use(ms...)
	g.GET("", card.GetAll)
}
