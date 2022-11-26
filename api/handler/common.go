package handler

import (
	b64 "encoding/base64"
	"encoding/json"
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
	// 1. Marshall
	stringJWT, err := json.Marshal(jwt)
	if err != nil {
		return err // TODO: use StringError(err)
	}

	// 2. base64 encode to avoid loosing bytes
	encodedJWT := b64.StdEncoding.EncodeToString(stringJWT)

	// 3. Create and set cookie
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = encodedJWT
	cookie.HttpOnly = true
	cookie.Expires = jwt.ExpAt
	c.SetCookie(cookie)

	return nil
}

func ReadJWTCookie(c echo.Context) (service.JWT, error) {
	jwt := service.JWT{}

	// 1. get cookie
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return jwt, err // TODO: use StringError(err)
	}

	// 2. decode
	decodedJWT, _ := b64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return jwt, err // TODO: use StringError(err)
	}

	// 3. unmarshal
	err = json.Unmarshal(decodedJWT, &jwt)
	if err != nil {
		return jwt, err // TODO: use StringError(err)
	}

	return jwt, nil
}

type HttpError struct {
	Error string `json:"error"`
}
