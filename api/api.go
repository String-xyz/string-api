package api

import (
	"net/http"

	"github.com/String-xyz/go-lib/v2/database"
	libmiddleware "github.com/String-xyz/go-lib/v2/middleware"
	"github.com/String-xyz/go-lib/v2/validator"

	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/api/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/String-xyz/string-api/pkg/service"
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

// @title String API
// @version 1.0
// @description String API for executing transactions and managing users

// @contact.name String API Support
// @contact.url http://string.xyz
// @contact.email support@stringxyz.com

// @host string-api.xyz
// @BasePath /
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
	transactRoute(services, e)
	quoteRoute(services, e)
	userRoute(services, e)
	loginRoute(services, e)
	verificationRoute(services, e)
	cardRoute(services, e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}

func StartInternal(config APIConfig) {
	e := echo.New()
	e.Validator = validator.New()
	baseMiddleware(config.Logger, e)
	e.GET("/heartbeat", heartbeat)

	// initialize routes - A route group only needs access to the services layer. It doesn't need access to the repos layer
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func baseMiddleware(logger *zerolog.Logger, e *echo.Echo) {
	e.Use(libmiddleware.Recover())
	e.Use(libmiddleware.RequestId())
	e.Use(libmiddleware.Tracer("api"))
	e.Use(libmiddleware.CORS())
	e.Use(libmiddleware.Logger(logger))
	e.Use(libmiddleware.LogRequest())
}

func transactRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewTransaction(e, services.Transaction)
	handler.RegisterRoutes(e.Group("/transactions"), middleware.JWTAuth())
}

func userRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewUser(e, services.User, services.Verification)
	handler.RegisterRoutes(e.Group("/users"), middleware.APIKeyPublicAuth(services.Auth), middleware.JWTAuth())
	handler.RegisterPrivateRoutes(e.Group("/users"), middleware.APIKeySecretAuth(services.Auth))
}

func loginRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewLogin(e, services.Auth, services.Device)
	handler.RegisterRoutes(e.Group("/login"), middleware.APIKeyPublicAuth(services.Auth))
}

func verificationRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewVerification(e, services.Verification, services.Device)
	handler.RegisterRoutes(e.Group("/verification"))
}

func quoteRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewQuote(e, services.Transaction)
	handler.RegisterRoutes(e.Group("/quotes"), middleware.JWTAuth())
}

func cardRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewCard(e, services.Card)
	handler.RegisterRoutes(e.Group("/cards"), middleware.JWTAuth())
}
