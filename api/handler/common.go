package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	service "github.com/String-xyz/string-api/pkg/service"
	"golang.org/x/crypto/sha3"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogError(c echo.Context, err error, handlerMsg string) {
	lg := c.Get("logger").(*zerolog.Logger)
	sp, _ := tracer.SpanFromContext(c.Request().Context())
	lg.Error().Stack().Err(err).Uint64("trace_id", sp.Context().TraceID()).
		Uint64("span_id", sp.Context().SpanID()).Msg(handlerMsg)
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

	if IsLocalEnv() {
		st2 := fmt.Sprintf("\nSTACK TRACE:\n%+v: [%+v ]\n\n", cause.Error(), st[0:5])
		// delete the string_api docker path from the stack trace
		st2 = strings.ReplaceAll(st2, "/string_api/", "")
		fmt.Print(st2)
		return
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
	cookie.Path = "/"             // Send cookie in every sub path request
	cookie.Secure = !IsLocalEnv() // in production allow https only
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
	cookie.Path = "/login/"       // Send cookie only in /login path request
	cookie.Secure = !IsLocalEnv() // in production allow https only
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

func DeleteAuthCookies(c echo.Context) error {
	// in order to delete a cookie we need to set it with an expired date
	cookie := new(http.Cookie)
	cookie.Name = "StringJWT"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/" // Send cookie in every sub path request
	cookie.Secure = !IsLocalEnv()
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/login/" // Send cookie only in refresh path request
	cookie.Secure = !IsLocalEnv()
	c.SetCookie(cookie)

	return nil
}

func IsLocalEnv() bool {
	return os.Getenv("ENV") == "local"
}

func validAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func SanitizeChecksums(addrs ...*string) {
	for _, addr := range addrs {
		if !validAddress(*addr) {
			continue
		}
		lowerCase := strings.ToLower(*addr)[2:]
		hash := sha3.NewLegacyKeccak256()
		hash.Write([]byte(lowerCase))
		hashBytes := hash.Sum(nil)

		valid := "0x"
		for i, b := range lowerCase {
			c := string(b)
			if b < '0' || b > '9' {
				if hashBytes[i/2]&byte(128-i%2*120) != 0 {
					c = string(b - 32)
				}
			}
			valid += c
		}
		*addr = valid
	}
}
