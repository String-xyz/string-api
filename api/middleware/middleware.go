package middleware

import (
	"net/http"
	"os"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func BearerAuth() echo.MiddlewareFunc {
	config := echoMiddleware.JWTConfig{
		TokenLookup: "header:Authorization,cookie:StringJWT",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			var claims = &service.JWTClaims{}
			t, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			c.Set("userId", claims.UserId)
			c.Set("deviceId", claims.DeviceId)
			c.Set("platformId", claims.PlatformId)
			return t, err
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		ErrorHandlerWithContext: func(err error, c echo.Context) error {

			return httperror.Unauthorized(c)
		},
	}
	return echoMiddleware.JWTWithConfig(config)
}

func APIKeyAuth(service service.Auth) echo.MiddlewareFunc {
	config := echoMiddleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(auth string, c echo.Context) (bool, error) {
			platformId, err := service.ValidateAPIKey(auth)
			if err != nil {
				return false, err
			}

			c.Set("platformId", platformId)

			return true, nil
		},
	}
	return echoMiddleware.KeyAuthWithConfig(config)
}

func Georestrict(service service.Geofencing) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			/* Get Ip from request */
			ip := c.RealIP()

			// check if location ip is restricted
			isAllowed, err := service.IsAllowed(ip)
			// in case of error, what we should do? allow or deny?
			// For now we are denying
			if err != nil || !isAllowed {
				if err != nil {
					libcommon.LogStringError(c, err, "Error in georestrict middleware")
				}
				return c.JSON(http.StatusForbidden, "Error: Geo Location Forbidden")
			}

			return next(c)
		}
	}
}
