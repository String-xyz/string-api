package api

import (
	"errors"

	"github.com/String-xyz/string-api/src/api/handler"
	"github.com/String-xyz/string-api/src/repository"
	"github.com/String-xyz/string-api/src/service"
	"github.com/labstack/echo/v4"
)

type APIConfig struct {
	DB   any
	Port string
	// redis
}

func Start(config APIConfig) error {
	e := echo.New()

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
