package handler

import (
	"fmt"
	"net/http"

	service "github.com/String-xyz/string-api/pkg/service"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogError(c echo.Context, err error, handlerMsg string) {
	lg := c.Get("logger").(*zerolog.Logger)
	lg.Err(err).Stack().Msg(handlerMsg)
}

func LogStringError(c echo.Context, err error, handlerMsg string) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	tracer, ok := errors.Cause(err).(stackTracer)
	if !ok {
		log.Warn().Str("error", err.Error()).Msg("error does not implement stack trace")
	}
	cause := errors.Cause(err)

	st := tracer.StackTrace()
	fmt.Printf("\n%+v: [%+v ]\n\n", cause.Error(), st[0:3])
	LogError(c, cause, handlerMsg)
}

func SetJWTCookie(c echo.Context, jwt service.JWT) error {
	cookie := new(http.Cookie)
	cookie.Name = "JWTCookie"
	cookie.Value = jwt.Token
	cookie.HttpOnly = true
	cookie.Expires = jwt.ExpAt
	// cookie.Secure = true
	c.SetCookie(cookie)

	return nil
}

type HttpError struct {
	Error string `json:"error"`
}
