package api

import (
	"net/http"

	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/api/middleware"
	"github.com/String-xyz/string-api/api/validator"
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
	e.Validator = validator.New()
	baseMiddleware(config.Logger, e)

	// not internal middlewares
	geofencingService := service.NewGeofencing(config.Redis)
	e.Use(middleware.Georestrict(geofencingService))

	e.GET("/heartbeat", heartbeat)
	repos := NewRepos(config)
	authService := service.NewAuth(repos)

	AuthAPIKey(config, e, true)
	transactRoute(config, repos, authService, e)
	quoteRoute(config, repos, authService, e)
	userRoute(repos, authService, e)
	loginRoute(repos, e)
	verificationRoute(repos, e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}

func StartInternal(config APIConfig) {
	e := echo.New()
	e.Validator = validator.New()
	baseMiddleware(config.Logger, e)
	e.GET("/heartbeat", heartbeat)
	platformRoute(config, e)
	AuthAPIKey(config, e, true)
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

func platformRoute(config APIConfig, e *echo.Echo) {
	p := repository.NewPlatform(config.DB)
	a := repository.NewAuth(config.Redis, config.DB)
	c := repository.NewContact(config.DB)
	service := service.NewPlatform(p, c, a)
	handler := handler.NewPlatform(service)
	handler.RegisterRoutes(e.Group("/platforms"), middleware.BearerAuth())
}

func AuthAPIKey(config APIConfig, e *echo.Echo, internal bool) {
	a := repository.NewAuth(config.Redis, config.DB)
	service := service.NewAPIKeyStrategy(a)
	handler := handler.NewAuthAPIKey(service, internal)
	handler.RegisterRoutes(e.Group("/apikeys"))
}

func transactRoute(config APIConfig, repos repository.Repositories, auth service.Auth, e *echo.Echo) {
	service := service.NewTransaction(repos, config.Redis)
	handler := handler.NewTransaction(e, service)
	handler.RegisterRoutes(e.Group("/transactions"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}

func userRoute(repos repository.Repositories, auth service.Auth, e *echo.Echo) {
	// user := service.NewUser(repos)
	user := NewServices(repos).User
	verification := service.NewVerification(repos)
	handler := handler.NewUser(e, user, verification)
	handler.RegisterRoutes(e.Group("/users"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}

func loginRoute(repos repository.Repositories, e *echo.Echo) {
	service := service.NewAuth(repos)
	handler := handler.NewLogin(e, service)
	handler.RegisterRoutes(e.Group("/login"))
}

func verificationRoute(repos repository.Repositories, e *echo.Echo) {
	verificationRepos := repository.Repositories{Contact: repos.Contact, User: repos.User}
	verification := service.NewVerification(verificationRepos)
	handler := handler.NewVerification(e, verification)
	handler.RegisterRoutes(e.Group("/verification"))
}

func quoteRoute(config APIConfig, repos repository.Repositories, auth service.Auth, e *echo.Echo) {
	service := service.NewTransaction(repos, config.Redis)
	handler := handler.NewQuote(e, service)
	handler.RegisterRoutes(e.Group("/quotes"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}
