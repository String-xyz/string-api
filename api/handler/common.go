package handler

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
	"golang.org/x/crypto/sha3"

	"github.com/labstack/echo/v4"
)

func SetJWTCookie(c echo.Context, jwt model.JWT) error {
	cookie := new(http.Cookie)
	cookie.Name = "StringJWT"
	cookie.Value = jwt.Token
	cookie.HttpOnly = true
	cookie.Expires = jwt.ExpAt // we want the cookie to expire at the same time as the token
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/"                    // Send cookie in every sub path request
	cookie.Secure = !common.IsLocalEnv() // in production allow https only
	c.SetCookie(cookie)

	return nil
}

func SetRefreshTokenCookie(c echo.Context, refresh model.RefreshTokenResponse) error {
	cookie := new(http.Cookie)
	cookie.Name = "StringRefreshToken"
	cookie.Value = refresh.Token
	cookie.HttpOnly = true
	cookie.Expires = refresh.ExpAt // we want the cookie to expire at the same time as the token
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/login/"              // Send cookie only in /login path request
	cookie.Secure = !common.IsLocalEnv() // in production allow https only
	c.SetCookie(cookie)

	return nil
}

func SetAuthCookies(c echo.Context, jwt model.JWT) error {
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
	cookie.HttpOnly = true
	cookie.Expires = time.Now()
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/" // Send cookie in every sub path request
	cookie.Secure = !common.IsLocalEnv()
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = "StringRefreshToken"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Expires = time.Now()
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/login/" // Send cookie only in refresh path request
	cookie.Secure = !common.IsLocalEnv()
	c.SetCookie(cookie)

	return nil
}

func validAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func getCookieSameSiteMode() http.SameSite {
	sameSiteMode := http.SameSiteNoneMode // allow cors
	if common.IsLocalEnv() {
		sameSiteMode = http.SameSiteLaxMode // because SameSiteNoneMode is not allowed in localhost we use lax mode
	}
	return sameSiteMode
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

func DefaultErrorHandler(c echo.Context, err error, handlerName string) error {
	if err == nil {
		return nil
	}

	// always log the error
	common.LogStringError(c, err, handlerName)

	if serror.Is(err, serror.NOT_FOUND) {
		return httperror.NotFound404(c)
	}

	if serror.Is(err, serror.FORBIDDEN) {
		return httperror.Forbidden403(c, "Invoking member lacks authority")
	}

	if serror.Is(err, serror.INVALID_RESET_TOKEN) {
		return httperror.BadRequest400(c, "Invalid password reset token")
	}

	if serror.Is(err, serror.INVALID_PASSWORD) {
		return httperror.BadRequest400(c, "Invalid password")
	}

	if serror.Is(err, serror.ALREADY_IN_USE) {
		return httperror.Conflict409(c, "Already in use")
	}

	if serror.Is(err, serror.INVALID_DATA) {
		return httperror.BadRequest400(c, "Invalid data")
	}

	return httperror.Internal500(c)
}
