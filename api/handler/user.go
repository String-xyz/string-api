package handler

import (
	b64 "encoding/base64"
	"net/http"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	serror "github.com/String-xyz/go-lib/stringerror"
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
	platformId := c.Get("platformId").(string)

	ctx := c.Request().Context()
	var body model.WalletSignaturePayloadSigned
	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user:create user bind")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	// base64 decode nonce
	decodedNonce, _ := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "user: create user decode nonce")
		return httperror.BadRequestError(c)
	}
	body.Nonce = string(decodedNonce)

	resp, err := u.userService.Create(ctx, body, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "user: creating user")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.ConflictError(c)
		}

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFoundError(c)
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.ForbiddenError(c, "Link expired, please request a new one")
		}

		return httperror.InternalError(c)
	}
	// set auth cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		libcommon.LogStringError(c, err, "user: unable to set auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (u user) Status(c echo.Context) error {
	ctx := c.Request().Context()
	valid, userId := validUserId(IdParam(c), c)
	if !valid {
		return httperror.Unauthorized(c)
	}

	status, err := u.userService.GetStatus(ctx, userId)
	if err != nil {
		libcommon.LogStringError(c, err, "user: get status")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, status)
}

func (u user) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var body model.UpdateUserName
	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: update bind")
		return httperror.BadRequestError(c)
	}
	_, userId := validUserId(IdParam(c), c)
	user, err := u.userService.Update(ctx, userId, body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: update")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, user)
}

// VerifyEmail send an email with a link, the user must click on the link for the email to be verified
// the link sent is handled by (verification.VerifyEmail) handler
func (u user) VerifyEmail(c echo.Context) error {
	ctx := c.Request().Context()
	platformId := c.Get("platformId").(string)
	_, userId := validUserId(IdParam(c), c)
	email := c.QueryParam("email")
	if email == "" {
		return httperror.BadRequestError(c, "Missing or invalid email")
	}

	err := u.verificationService.SendEmailVerification(ctx, userId, email, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "user: email verification")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.ConflictError(c)
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.ForbiddenError(c, "Link expired, please request a new one")
		}

		return httperror.InternalError(c, "Unable to send email verification")
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
	g.POST("", u.Create, ms[0])

	// the rest of the endpoints do not require api key
	ms = ms[1:]

	g.GET("/:id/status", u.Status, ms...)
	g.GET("/:id/verify-email", u.VerifyEmail, ms...)
	g.PUT("/:id", u.Update, ms...)
}

// get userId from context and also compare if both are valid
// this is useful for path params validation and userId from JWT
func validUserId(userId string, c echo.Context) (bool, string) {
	_userId := c.Get("userId").(string)
	return _userId == userId, _userId
}

func IdParam(c echo.Context) string {
	return c.Param("id")
}
