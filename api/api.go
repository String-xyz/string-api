package api

import (
	"errors"
	"net/http"

	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/api/middleware"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type APIConfig struct {
	DB   *sqlx.DB
	Port string
	// redis
}

func heartbeat(c echo.Context) error {
	return c.JSON(http.StatusOK, "healthy")
}

func Start(config APIConfig) error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/heartbeat", heartbeat)
	// Allow all CORS

	// Todo: Add Logger

	// TODO: Create middleware for jwt

	////////////////////////
	// REPOSITORIES
	////////////////////////
	transactRepo := repository.NewTransaction(config.DB)

	////////////////////////
	// SERVICES
	////////////////////////
	transactService := service.NewTransaction(transactRepo)

	////////////////////////
	// HANDLERS
	////////////////////////
	transactHandler := handler.NewTransaction(e, transactService)

	////////////////////////
	// MIDDLEWARE
	////////////////////////
	// cors := middleware.StringMiddleware(middleware.MiddlewareConfig{any: 0})

	var transactMiddleware []echo.MiddlewareFunc

	////////////////////////
	// REGISTER ROUTES
	////////////////////////
	transactHandler.RegisterRoutes(e.Group("/transact"), transactMiddleware...)

	e.Logger.Fatal(e.Start(":" + config.Port))

	return errors.New("error")

}
