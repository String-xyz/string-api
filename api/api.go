package api

import (
	"net/http"

	libCommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/middleware"
	validator "github.com/String-xyz/go-lib/validator"
	"github.com/String-xyz/string-api/api/handler"
	libMiddleware "github.com/String-xyz/string-api/api/middleware"

	"github.com/String-xyz/string-api/pkg/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type APIConfig struct {
	DB     *sqlx.DB
	Redis  database.RedisStore
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
	e.Use(libMiddleware.Georestrict(geofencingService))

	e.GET("/heartbeat", heartbeat)

	// initialize route dependencies
	repos := NewRepos(config)
	services := NewServices(config, repos)

	// initialize routes - A route group only needs access to the services layer. It should'n access the repos layer directly
	AuthAPIKey(services, e, libCommon.IsLocalEnv())
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
	e.Use(middleware.Tracer())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger(logger))
	e.Use(middleware.LogRequest())
}

func platformRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewPlatform(services.Platform)
	handler.RegisterRoutes(e.Group("/platforms"), libMiddleware.BearerAuth())
}

func AuthAPIKey(services service.Services, e *echo.Echo, internal bool) {
	handler := handler.NewAuthAPIKey(services.ApiKey, internal)
	handler.RegisterRoutes(e.Group("/apikeys"))
}

func transactRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewTransaction(e, services.Transaction)
	handler.RegisterRoutes(e.Group("/transactions"), libMiddleware.APIKeyAuth(services.Auth), libMiddleware.BearerAuth())
}

func userRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewUser(e, services.User, services.Verification)
	handler.RegisterRoutes(e.Group("/users"), libMiddleware.APIKeyAuth(services.Auth), libMiddleware.BearerAuth())
}

func loginRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewLogin(e, services.Auth, services.Device)
	handler.RegisterRoutes(e.Group("/login"), libMiddleware.APIKeyAuth(services.Auth))
}

func verificationRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewVerification(e, services.Verification, services.Device)
	handler.RegisterRoutes(e.Group("/verification"))
}

func quoteRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewQuote(e, services.Transaction)
	handler.RegisterRoutes(e.Group("/quotes"), libMiddleware.APIKeyAuth(services.Auth), libMiddleware.BearerAuth())
}
