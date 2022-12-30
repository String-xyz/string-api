package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	service "github.com/String-xyz/string-api/pkg/service"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogError(c echo.Context, err error, handlerMsg string) {
	lg := c.Get("logger").(*zerolog.Logger)
	lg.Error().Stack().Err(err).Msg(handlerMsg)
}

func LogStringError(c echo.Context, err error, handlerMsg string) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	tracer, ok := errors.Cause(err).(stackTracer)
	if !ok {
		log.Warn().Str("error", err.Error()).Msg("error does not implement stack trace")
		return
	}

	cause := errors.Cause(err)
	st := tracer.StackTrace()

	if os.Getenv("env") == "local" {
		st2 := fmt.Sprintf("\n%+v: [%+v ]\n\n", cause.Error(), st[0:3])
		// delete the string_api docker path from the stack trace
		st2 = strings.ReplaceAll(st2, "/string_api/", "")
		fmt.Print(st2)
	}

	LogError(c, err, handlerMsg)
}

func SetJWTCookie(c echo.Context, jwt service.JWT) error {
	cookie := new(http.Cookie)
	cookie.Name = "StringJWT"
	cookie.Value = jwt.Token
	// cookie.HttpOnly = true // due the short expiration time it is not needed to be http only
	cookie.Expires = jwt.ExpAt // we want the cookie to expire at the same time as the token
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"              // Send cookie in every sub path request
	cookie.Secure = isProduction() // in production allow https only
	c.SetCookie(cookie)

	return nil
}

func SetRefreshTokenCookie(c echo.Context, refresh service.RefreshTokenResponse) error {
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refresh.Token
	cookie.HttpOnly = true
	cookie.Expires = refresh.ExpAt // we want the cookie to expire at the same time as the token
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/login/refresh" // Send cookie only in refresh path request
	cookie.Secure = isProduction() // in production allow https only
	c.SetCookie(cookie)

	return nil
}

func SetAuthCookies(c echo.Context, jwt service.JWT) error {
	err := SetJWTCookie(c, jwt)
	if err != nil {
		return err
	}

	err = SetRefreshTokenCookie(c, jwt.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}

func isProduction() bool {
	return os.Getenv("ENV") == "production"
}
