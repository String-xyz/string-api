package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type User interface {
	GetStatus(c echo.Context) error    // If wallet addr is associated with user, return current state of their onboarding
	Create(c echo.Context) error       // Create new user using wallet addr, optionally mark as validated if signature is provided
	Sign(c echo.Context) error         // Takes in a signed timestamp from user, validating their wallet
	Authenticate(c echo.Context) error // Takes e-mail and wallet addr of user, validates email with twilio
	Name(c echo.Context) error         // Takes name and wallet addr of user, associates name with wallet addr
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
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
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, res)
}

func (u user) Create(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = u.Service.Create(body)
	if err != nil {
		lg.Err(err).Msg("user create")
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, nil)
}

func (u user) Sign(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = u.Service.Sign(body)
	if err != nil {
		lg.Err(err).Msg("user sign")
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, nil)
}

func (u user) Authenticate(c echo.Context) error {
	lg := c.Get("logger").(*zerolog.Logger)
	var body model.UserRequest
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	// Token was provided
	token := c.QueryParam("token")
	if token != "" {
		err = u.Service.ReceiveEmailAuthentication(token)
		if err != nil {
			lg.Err(err).Msg("user authenticate")
			return c.String(http.StatusOK, "User Service Failed")
		}
		return c.JSON(http.StatusOK, nil)
	}

	// User needs a token
	err = u.Service.Authenticate(body)
	if err != nil {
		lg.Err(err).Msg("user authenticate")
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, nil)
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
		return c.String(http.StatusOK, "User Service Failed")
	}
	return c.JSON(http.StatusOK, nil)
}

func (u user) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	u.Group = g
	g.Use(ms...)
	g.GET("/", u.GetStatus)
	g.POST("/", u.Create)
	g.PUT("/", u.Sign)
	g.POST("/email", u.Authenticate)
	g.POST("/name", u.Name)
}
