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
)

type APIConfig struct {
	DB    *sqlx.DB
	Redis store.RedisStore
	Port  string
}

func heartbeat(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

func Start(config APIConfig) {
	e := echo.New()
	defaultMiddleware(e)
	e.GET("/heartbeat", heartbeat)
	transactRepo := repository.NewTransaction(config.DB)
	transactService := service.NewTransaction(transactRepo)
	transactHandler := handler.NewTransaction(e, transactService)
	transactHandler.RegisterRoutes(e.Group("/transact"), middleware.Auth())
	authRoute(config, e)
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func defaultMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
}

func authRoute(config APIConfig, e *echo.Echo) {
	auth := repository.NewAuth(config.Redis, config.DB)
	user := repository.NewUser(config.DB)
	contact := repository.NewUserContact(config.DB)
	service := service.NewAuth(auth, user, contact)
	handler := handler.NewAuth(service)
	handler.RegisterRoutes(e.Group("/auth"))
}
