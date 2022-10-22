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
	Logger *zerolog.Logger
	Port   string
}

func heartbeat(c echo.Context) error {
	return c.JSON(http.StatusOK, "alive")
}

func Start(config APIConfig) {
	e := echo.New()
	baseMiddleware(config.Logger, e)
	e.GET("/heartbeat", heartbeat)
	authService := authRoute(config, e)
	platformRoute(config, e)
	transactRoute(config, authService, e)
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func baseMiddleware(logger *zerolog.Logger, e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Tracer())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger(logger))
	e.Use(middleware.LogRequest())
}

func authRoute(config APIConfig, e *echo.Echo) service.Auth {
	a := repository.NewAuth(config.Redis, config.DB)
	u := repository.NewUser(config.DB)
	c := repository.NewUserContact(config.DB)
	service := service.NewAuth(a, u, c)
	handler := handler.NewAuth(service)
	handler.RegisterRoutes(e.Group("/auth"))
	return service
}

func platformRoute(config APIConfig, e *echo.Echo) {
	p := repository.NewPlatform(config.DB)
	a := repository.NewAuth(config.Redis, config.DB)
	c := repository.NewUserContact(config.DB)
	service := service.NewPlatform(p, c, a)
	handler := handler.NewPlatform(service)
	handler.RegisterRoutes(e.Group("/platform"), middleware.BearerAuth())
}

func transactRoute(config APIConfig, auth service.Auth, e *echo.Echo) {
	repos := service.TransactionRepos{
		Asset:       repository.NewAsset(config.DB),
		Network:     repository.NewNetwork(config.DB),
		Transaction: repository.NewTransaction(config.DB),
		TxLeg:       repository.NewTxLeg(config.DB),
		User:        repository.NewUser(config.DB),
		// More will follow
	}
	service := service.NewTransaction(repos)
	handler := handler.NewTransaction(e, service)
	handler.RegisterRoutes(e.Group("/transact"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}
