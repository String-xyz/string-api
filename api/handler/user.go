package handler

import (
	"net/http"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type User interface {
	Create(c echo.Context) error
	Status(c echo.Context) error
	Update(c echo.Context) error
	VerifyEmail(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type ResultMessage struct {
	Status string
}

type user struct {
	userService         service.User
	verificationService service.Verification
	Group               *echo.Group
}

func NewUser(route *echo.Echo, userSrv service.User, verificationSrv service.Verification) User {
	return &user{userSrv, verificationSrv, nil}
}

func (u user) Create(c echo.Context) error {
	var body model.WalletSignaturePayload
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "user:create user bind")
		return c.JSON(http.StatusBadRequest, HttpError{Error: "Missing or invalid body"})
	}

	jwt, err := u.userService.Create(body)
	if err != nil {
		LogStringError(c, err, "user: creating user")
		return c.JSON(http.StatusInternalServerError, HttpError{Error: "Error creating user"})
	}
	return c.JSON(http.StatusOK, jwt)
}

func (u user) Status(c echo.Context) error {
	valid, userId := validUserID(IDParam(c), c)
	if !valid {
		return c.JSON(http.StatusUnauthorized, HttpError{Error: "Unauthorized"})
	}
	walletAddress := c.QueryParam("walletAddress")
	if walletAddress == "" {
		return c.JSON(http.StatusBadRequest, HttpError{Error: "Missing or invalid walletAddress"})
	}
	status, err := u.userService.GetStatus(userId, walletAddress)
	if err != nil {
		LogStringError(c, err, "user: get status")
		return c.JSON(http.StatusInternalServerError, HttpError{Error: "Error getting status"})
	}
	return c.JSON(http.StatusOK, status)
}

func (u user) Update(c echo.Context) error {
	var body model.UpdateUserName
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "user: update bind")
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = u.userService.Update(body)
	if err != nil {
		LogStringError(c, err, "user: update")
		return c.String(http.StatusBadRequest, "Error updating user")
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "User updated successfully"})
}

// VerifyEmail send an email with a link, the user must click on the link for the email to be verified
// the link sent is handled by (verification.VerifyEmail) handler
func (u user) VerifyEmail(c echo.Context) error {
	_, userId := validUserID(IDParam(c), c)
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, HttpError{Error: "Missing or invalid email"})
	}
	err := u.verificationService.SendEmailVerification(userId, email)
	if err != nil {
		LogStringError(c, err, "user: email verification")
		return c.JSON(http.StatusInternalServerError, HttpError{Error: "Could not send email verification"})
	}
	return c.JSON(http.StatusOK, ResultMessage{Status: "Email verification sent"})
}

func (u user) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	u.Group = g
	// create does not require JWT auth middleware
	// hence the placing before adding the apiKey middleware
	if len(ms) >= 2 {
		g.Use(ms[0])
		g.POST("", u.Create)
		g.Use(ms[1])
	} else {
		g.Use(ms...)
		g.POST("", u.Create)
	}
	g.GET("/:id/status", u.Status)
	g.GET("/:id/verify-email", u.VerifyEmail)
	g.PUT("/:id", u.Update)
}

// get userId from context and also compare if both are valid
// this is useful for path params validation and userID from JWT
func validUserID(userID string, c echo.Context) (bool, string) {
	userId := c.Get("userId").(string)
	return userId == userID, userId
}

func IDParam(c echo.Context) string {
	return c.Param("id")
}
