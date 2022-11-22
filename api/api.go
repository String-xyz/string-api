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

	// not internal middlewares
	geofencingService := service.NewGeofencing(config.Redis)
	e.Use(middleware.Georestrict(geofencingService))

	e.GET("/heartbeat", heartbeat)
	authService := authRoute(config, e)
	AuthAPIKey(config, e, true)
	transactRoute(config, authService, e)
	userRoute(config, authService, e)

	loginRoute(config, e)
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func StartInternal(config APIConfig) {
	e := echo.New()
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

func authRoute(config APIConfig, e *echo.Echo) service.Auth {
	a := repository.NewAuth(config.Redis, config.DB)
	service := service.NewAuth(a)
	handler := handler.NewAuth(service)
	handler.RegisterRoutes(e.Group("/auth"))
	return service
}

func platformRoute(config APIConfig, e *echo.Echo) {
	p := repository.NewPlatform(config.DB)
	a := repository.NewAuth(config.Redis, config.DB)
	c := repository.NewContact(config.DB)
	service := service.NewPlatform(p, c, a)
	handler := handler.NewPlatform(service)
	handler.RegisterRoutes(e.Group("/platform"), middleware.BearerAuth())
}

func AuthAPIKey(config APIConfig, e *echo.Echo, internal bool) {
	a := repository.NewAuth(config.Redis, config.DB)
	service := service.NewAPIKeyStrategy(a)
	handler := handler.NewAuthAPIKey(service, internal)
	handler.RegisterRoutes(e.Group("/apikey"))
}

func transactRoute(config APIConfig, auth service.Auth, e *echo.Echo) {
	repos := service.TransactionRepos{
		Asset:       repository.NewAsset(config.DB),
		Network:     repository.NewNetwork(config.DB),
		Transaction: repository.NewTransaction(config.DB),
		TxLeg:       repository.NewTxLeg(config.DB),
		User:        repository.NewUser(config.DB),
		Instrument:  repository.NewInstrument(config.DB),
		Device:      repository.NewDevice(config.DB),
		Location:    repository.NewLocation(config.DB),
	}
	service := service.NewTransaction(repos)
	handler := handler.NewTransaction(e, service)
	handler.RegisterRoutes(e.Group("/transact"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}

func userRoute(config APIConfig, auth service.Auth, e *echo.Echo) {
	repos := service.UserRepos{
		Auth:         repository.NewAuth(config.Redis, config.DB),
		User:         repository.NewUser(config.DB),
		Contact:      repository.NewContact(config.DB),
		Instrument:   repository.NewInstrument(config.DB),
		Device:       repository.NewDevice(config.DB),
		UserPlatform: repository.NewUserPlatform(config.DB),
	}
	service := service.NewUser(repos)
	handler := handler.NewUser(e, service)
	handler.RegisterRoutes(e.Group("/user"), middleware.APIKeyAuth(auth), middleware.BearerAuth())
}

func loginRoute(config APIConfig, e *echo.Echo) {
	repos := service.UserRepos{
		Auth:         repository.NewAuth(config.Redis, config.DB),
		User:         repository.NewUser(config.DB),
		Contact:      repository.NewContact(config.DB),
		Instrument:   repository.NewInstrument(config.DB),
		Device:       repository.NewDevice(config.DB),
		UserPlatform: repository.NewUserPlatform(config.DB),
	}
	service := service.NewUser(repos)
	handler := handler.NewLogin(e, service)
	handler.RegisterRoutes(e.Group("/login"))
}
