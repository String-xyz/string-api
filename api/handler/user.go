package handler

import (
	b64 "encoding/base64"
	"net/http"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/go-lib/v2/validator"
	"github.com/labstack/echo/v4"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
)

type User interface {
	Create(c echo.Context) error
	Status(c echo.Context) error
	Update(c echo.Context) error
	PreviewEmail(c echo.Context) error
	VerifyEmail(c echo.Context) error
	PreValidateEmail(c echo.Context) error
	GetPersonaAccountId(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
	RegisterPrivateRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type ResultMessage struct {
	Status string
}

type user struct {
	userService  service.User
	verification service.Verification
	Group        *echo.Group
}

func NewUser(route *echo.Echo, userSrv service.User, verificationSrv service.Verification) User {
	return &user{userSrv, verificationSrv, nil}
}

// @Summary Create user
// @Description Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param body body model.WalletSignaturePayloadSigned true "Wallet Signature Payload Signed"
// @Success 200 {object} model.UserLoginResponse
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Router /users [post]
func (u user) Create(c echo.Context) error {
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	ctx := c.Request().Context()
	var body model.WalletSignaturePayloadSigned
	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user:create user bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		libcommon.LogStringError(c, err, "user:create user validate body")
		return httperror.InvalidPayload400(c, err)
	}

	// base64 decode nonce
	decodedNonce, _ := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "user: create user decode nonce")
		return httperror.BadRequest400(c)
	}
	body.Nonce = string(decodedNonce)

	resp, err := u.userService.Create(ctx, body, platformId)
	if err != nil {
		libcommon.LogStringError(c, err, "user: creating user")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.Conflict409(c)
		}

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c)
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.Forbidden403(c, "Nonce expired. Request a new one")
		}

		return httperror.Internal500(c)
	}
	// set auth cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		libcommon.LogStringError(c, err, "user: unable to set auth cookies")
		return httperror.Internal500(c)
	}

	// 200
	return c.JSON(http.StatusOK, resp)
}

// @Summary Get user status
// @Description Get user status
// @Tags Users
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "User ID"
// @Success 200 {object} model.UserOnboardingStatus
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /users/{id}/status [get]
func (u user) Status(c echo.Context) error {
	ctx := c.Request().Context()
	valid, userId := validUserId(IdParam(c), c)
	if !valid {
		return httperror.Unauthorized401(c)
	}

	status, err := u.userService.GetStatus(ctx, userId)
	if err != nil {
		libcommon.LogStringError(c, err, "user: get status")
		return httperror.Internal500(c)
	}
	// 200
	return c.JSON(http.StatusOK, status)
}

// @Summary Update user
// @Description Update user
// @Tags Users
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "User ID"
// @Param body body model.UpdateUserName true "Update User Name"
// @Success 200 {object} model.User
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /users/{id} [patch]
func (u user) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var body model.UpdateUserName
	platformId, ok := c.Get("platformId").(string)

	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: update bind")
		return httperror.BadRequest400(c)
	}

	err = c.Validate(body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: update validate body")
		return httperror.InvalidPayload400(c, err)
	}

	_, userId := validUserId(IdParam(c), c)

	user, err := u.userService.Update(ctx, userId, platformId, body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: update")
		return httperror.Internal500(c)
	}

	// 200
	return c.JSON(http.StatusOK, user)
}

// @Summary Verify email
// @Description Verify email sends an email with a link, the user must click on the link for the email to be verified. The link sent is handled by (verification.VerifyEmail) handler
// @Tags Users
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "User ID"
// @Param email query string true "Email to verify"
// @Success 200 {object} ResultMessage
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Router /users/{id}/verify-email [get]
func (u user) VerifyEmail(c echo.Context) error {
	ctx := c.Request().Context()
	email := c.QueryParam("email")
	platformId, ok := c.Get("platformId").(string)

	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	valid, userId := validUserId(IdParam(c), c)
	if !valid {
		return httperror.BadRequest400(c, "Missing or invalid user id")
	}

	if !validator.ValidEmail(email) {
		return httperror.BadRequest400(c, "Invalid email")
	}

	err := u.verification.SendEmailVerification(ctx, platformId, userId, email)
	if err != nil {
		libcommon.LogStringError(c, err, "user: email verification")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.Conflict409(c)
		}

		return httperror.Internal500(c, "Unable to send email verification")
	}

	// 200
	return c.JSON(http.StatusOK, ResultMessage{Status: "Verification email sent"})
}

// @Summary Request device verification
// @Description Sends an email with a link, the user must click on the link for the device to be verified.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
func (u user) RequestDeviceVerification(c echo.Context) error {
	ctx := c.Request().Context()

	var body model.WalletSignaturePayloadSigned

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: request device verify bind")
		return httperror.BadRequest400(c)
	}

	err = c.Validate(body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: request device verify validate body")
		return httperror.InvalidPayload400(c, err)
	}

	// base64 decode nonce
	decodedNonce, err := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "login: verify signature decode nonce")
		return httperror.BadRequest400(c)
	}
	body.Nonce = string(decodedNonce)

	err = u.userService.RequestDeviceVerification(ctx, body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: device verification")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c)
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.BadRequest400(c, "Expired, request a new payload")
		}

		return httperror.Internal500(c, "Unable to send device verification")
	}

	// 200
	return c.JSON(http.StatusOK, ResultMessage{Status: "Device verification email sent"})
}

func (u user) GetDeviceStatus(c echo.Context) error {
	ctx := c.Request().Context()

	var body model.WalletSignaturePayloadSigned

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: request device verify bind")
		return httperror.BadRequest400(c)
	}

	err = c.Validate(body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: request device verify validate body")
		return httperror.InvalidPayload400(c, err)
	}

	// base64 decode nonce
	decodedNonce, err := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "login: verify signature decode nonce")
		return httperror.BadRequest400(c)
	}
	body.Nonce = string(decodedNonce)

	status, err := u.userService.GetDeviceStatus(ctx, body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: get device status")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c)
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.BadRequest400(c, "Expired, request a new payload")
		}

		return httperror.BadRequest400(c, "Invalid Payload")
	}

	// 200
	return c.JSON(http.StatusOK, status)
}

// @Summary Pre validate email
// @Description Pre validate email allows an organization to pre validate an email before the user signs up
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param email query string true "Pre Validated Email"
// @Success 200 {object} ResultMessage
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Router /users/{id}/email/pre-validate [post]
func (u user) PreValidateEmail(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Param("id")
	email := c.QueryParam("email")
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid platformId")
	}

	if !validator.ValidEmail(email) {
		return httperror.BadRequest400(c, "Invalid email")
	}

	err := u.verification.PreValidateEmail(ctx, platformId, userId, email)
	if err != nil {
		// ?
		libcommon.LogStringError(c, err, "user: pre validate email")
		return DefaultErrorHandler(c, err, "platformInternal: PreValidateEmail")
	}

	// 200
	return c.JSON(http.StatusOK, ResultMessage{Status: "validated"})
}

// @Summary Get user persona account id
// @Description Get user persona account id
// @Tags Users
// @Accept json
// @Produce json
// @Security JWT
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /users/persona-account-id [get]
func (u user) GetPersonaAccountId(c echo.Context) error {
	ctx := c.Request().Context()

	userId, ok := c.Get("userId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid userId")
	}

	accountId, err := u.userService.GetPersonaAccountId(ctx, userId)
	if err != nil {
		libcommon.LogStringError(c, err, "user: get persona account id")
		return httperror.Internal500(c)
	}
	return c.JSON(http.StatusOK, accountId)
}

// @Summary Get user email preview
// @Description Get obscured user email
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
func (u user) PreviewEmail(c echo.Context) error {
	ctx := c.Request().Context()

	var body model.WalletSignaturePayloadSigned

	err := c.Bind(&body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: preview email bind")
		return httperror.BadRequest400(c)
	}

	err = c.Validate(body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: preview email validate body")
		return httperror.InvalidPayload400(c, err)
	}

	// base64 decode nonce
	decodedNonce, err := b64.URLEncoding.DecodeString(body.Nonce)
	if err != nil {
		libcommon.LogStringError(c, err, "login: verify signature decode nonce")
		return httperror.BadRequest400(c)
	}
	body.Nonce = string(decodedNonce)

	email, err := u.userService.PreviewEmail(ctx, body)
	if err != nil {
		libcommon.LogStringError(c, err, "user: preview email")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c)
		}

		if serror.Is(err, serror.EXPIRED) {
			return httperror.BadRequest400(c, "Expired, request a new payload")
		}

		return httperror.BadRequest400(c, "Invalid Payload")
	}

	// 200
	return c.JSON(http.StatusOK, email)
}

func (u user) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No group attached to the User Handler")
	}
	u.Group = g
	// These endpoints use the API key and do not require JWT auth middleware
	// hence adding only the first middleware only which is APIKey
	g.POST("", u.Create, ms[0])
	g.POST("/preview-email", u.PreviewEmail, ms[0])
	g.POST("/verify-device", u.RequestDeviceVerification, ms[0])
	g.POST("/device-status", u.GetDeviceStatus, ms[0])

	// the rest of the endpoints use the JWT auth and do not require an API Key
	// hence removing the first (API key) middleware
	ms = ms[1:]

	g.GET("/:id/status", u.Status, ms...)
	g.GET("/:id/verify-email", u.VerifyEmail, ms...)
	g.GET("/persona-account-id", u.GetPersonaAccountId, ms...)
	g.PATCH("/:id", u.Update, ms...)
}

func (u user) RegisterPrivateRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("No private group attached to the User Handler")
	}
	u.Group = g

	g.POST("/:id/email/pre-validate", u.PreValidateEmail, ms...)
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
