package api

import (
	"net/http"

	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/api/middleware"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type APIConfig struct {
	DB     *sqlx.DB
	Redis  store.RedisStore
	Logger zerolog.Logger
	Port   string
}

func heartbeat(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

func Start(config APIConfig) {
	e := echo.New()
	baseMiddleware(config.Logger, e)
	e.GET("/heartbeat", heartbeat)
	authService := authRoute(config, e)
	transactRoute(config, authService, e)
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func baseMiddleware(logger zerolog.Logger, e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger(logger))
}

func authRoute(config APIConfig, e *echo.Echo) service.Auth {
	auth := repository.NewAuth(config.Redis, config.DB)
	user := repository.NewUser(config.DB)
	contact := repository.NewUserContact(config.DB)
	service := service.NewAuth(auth, user, contact)
	handler := handler.NewAuth(service)
	handler.RegisterRoutes(e.Group("/auth"))
	return service
}

func transactRoute(config APIConfig, auth service.Auth, e *echo.Echo) {
	repo := repository.NewTransaction(config.DB)
	service := service.NewTransaction(repo)
	handler := handler.NewTransaction(e, service)
	handler.RegisterRoutes(e.Group("/transact"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}
