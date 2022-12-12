package api

import (
	"net/http"

	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/api/middleware"
	"github.com/String-xyz/string-api/api/validator"
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
	e.Validator = validator.New()
	baseMiddleware(config.Logger, e)

	// not internal middlewares
	geofencingService := service.NewGeofencing(config.Redis)
	e.Use(middleware.Georestrict(geofencingService))

	e.GET("/heartbeat", heartbeat)

	// initialize route dependencies
	repos := NewRepos(config)
	services := NewServices(config, repos)

	// initialize routes - A route group only needs access to the services layer. It should'n access the repos layer directly
	AuthAPIKey(services, e, true)
	transactRoute(services, e)
	quoteRoute(services, e)
	userRoute(services, e)
	loginRoute(services, e)
	verificationRoute(services, e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}

func StartInternal(config APIConfig) {
	e := echo.New()
	e.Validator = validator.New()
	baseMiddleware(config.Logger, e)
	e.GET("/heartbeat", heartbeat)

	// initialize route dependencies
	repos := NewRepos(config)
	services := NewServices(config, repos)

	// initialize routes - A route group only needs access to the services layer. It doesn't need access to the repos layer
	platformRoute(services, e)
	AuthAPIKey(services, e, true)
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func baseMiddleware(logger *zerolog.Logger, e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Tracer())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger(logger))
	e.Use(middleware.LogRequest())
}

func platformRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewPlatform(services.Platform)
	handler.RegisterRoutes(e.Group("/platforms"), middleware.BearerAuth())
}

func AuthAPIKey(services service.Services, e *echo.Echo, internal bool) {
	handler := handler.NewAuthAPIKey(services.ApiKey, internal)
	handler.RegisterRoutes(e.Group("/apikeys"))
}

func transactRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewTransaction(e, services.Transaction)
	handler.RegisterRoutes(e.Group("/transactions"), middleware.APIKeyAuth(services.Auth), middleware.BearerAuth())
}

func userRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewUser(e, services.User, services.Verification)
	handler.RegisterRoutes(e.Group("/users"), middleware.APIKeyAuth(services.Auth), middleware.BearerAuth())
}

func loginRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewLogin(e, services.Auth)
	handler.RegisterRoutes(e.Group("/login"))
}

func verificationRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewVerification(e, services.Verification)
	handler.RegisterRoutes(e.Group("/verification"))
}

func quoteRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewQuote(e, services.Transaction)
	handler.RegisterRoutes(e.Group("/quotes"), middleware.APIKeyAuth(services.Auth), middleware.BearerAuth())
}
