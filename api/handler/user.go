package handler

import (
	"net/http"
	"strings"

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
	var body model.WalletSignaturePayloadSigned
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "user:create user bind")
		return BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return InvalidPayloadError(c, err)
	}

	resp, err := u.userService.Create(body)
	if err != nil {
		if errors.Cause(err) == "wallet already associated with user" {
			return Conflict(c)
		}

		LogStringError(c, err, "user: creating user")
		return InternalError(c)
	}
	// set jwt in cookie
	err = SetJWTCookie(c, resp.JWT)
	if err != nil {
		LogStringError(c, err, "user: unable to set JWT cookie")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (u user) Status(c echo.Context) error {
	valid, userId := validUserID(IDParam(c), c)
	if !valid {
		return Unauthorized(c)
	}

	status, err := u.userService.GetStatus(userId)
	if err != nil {
		LogStringError(c, err, "user: get status")
		return InternalError(c)
	}
	return c.JSON(http.StatusOK, status)
}

func (u user) Update(c echo.Context) error {
	var body model.UpdateUserName
	err := c.Bind(&body)
	if err != nil {
		LogStringError(c, err, "user: update bind")
		return BadRequestError(c)
	}
	_, userId := validUserID(IDParam(c), c)
	user, err := u.userService.Update(userId, body)
	if err != nil {
		LogStringError(c, err, "user: update")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, user)
}

// VerifyEmail send an email with a link, the user must click on the link for the email to be verified
// the link sent is handled by (verification.VerifyEmail) handler
func (u user) VerifyEmail(c echo.Context) error {
	_, userId := validUserID(IDParam(c), c)
	email := c.QueryParam("email")
	if email == "" {
		return BadRequestError(c, "Missing or invalid email")
	}

	err := u.verificationService.SendEmailVerification(userId, email)
	if err != nil {
		LogStringError(c, err, "user: email verification")
		return InternalError(c, "Unable to send email verification")
	}

	return c.JSON(http.StatusOK, ResultMessage{Status: "Email Successfully Verified"})
}

func (u user) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	u.Group = g
	// create does not require JWT auth middleware
	// hence adding only the first middleware only which is APIKey
	// Please pass the middle in order of API -> JWT otherwise this wont work.
	if len(ms) > 0 {
		g.POST("", u.Create, ms[0])
	}
	g.GET("/:id/status", u.Status, ms...)
	g.GET("/:id/verify-email", u.VerifyEmail, ms...)
	g.PUT("/:id", u.Update, ms...)
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
