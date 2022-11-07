package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type User interface {
	GetStatus(c echo.Context) error // If wallet addr is associated with user, return current state of their onboarding
	Name(c echo.Context) error      // Takes name and wallet addr of user, associates name with wallet addr
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type ResultMessage struct {
	Status string
}

type user struct {
	Service service.User
	Group   *echo.Group
}

func NewUser(route *echo.Echo, service service.User) User {
	return &user{service, nil}
}

func (u user) GetStatus(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	res, err := u.Service.GetStatus(body)
	if err != nil {
		lg.Err(err).Msg("user getstatus")
		return c.String(http.StatusNotFound, "User Not Found")
	}
	return c.JSON(http.StatusOK, res)
}

func (u user) Name(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = u.Service.Name(body)
	if err != nil {
		lg.Err(err).Msg("user name")
		return c.String(http.StatusBadRequest, "Could Not Update Name")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Name Updated Successfully"})
}

func (u user) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	u.Group = g
	g.Use(ms...)
	g.GET("", u.GetStatus)
	g.POST("/name", u.Name)
}
