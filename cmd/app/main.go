package main

import (
	"log"
	"os"
	"regexp"

	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/repository"
	"github.com/String-xyz/string-api/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func allowOrigin(origin string) (bool, error) {
	// TODO: Modify to be more restrictive
	return regexp.MatchString(`*`, origin)
}

func main() {
	// TODO: create db connection to Postgres
	// Checkout https://github.com/jackc/pgx
	// Also GORM? https://gorm.io/index.html

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	e := echo.New()

	////////////////////////
	// MIDDLEWARE
	////////////////////////

	// Allow all CORS
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOriginFunc: allowOrigin,
	// 	AllowMethods:    []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	// }))

	// Todo: Add Logger

	// TODO: Create middleware for jwt

	transactRepo := repository.NewTransaction()
	transactService := service.NewTransaction(transactRepo)
	transactHandler := handler.NewTransactionHandler(e, transactService)
	transactHandler.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}
	e.Logger.Fatal(e.Start(":" + port))
}
