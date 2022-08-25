package main

import (
	"github.com/String-xyz/string-api/api/handler"
	"github.com/String-xyz/string-api/repository"
	"github.com/String-xyz/string-api/service"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	transactRepo := repository.NewTransaction()
	transactService := service.NewTransactor(transactRepo)
	transactHandler := handler.NewTransactionHandler(e, transactService)
	transactHandler.RegisterRoutes()

	e.Logger.Fatal(e.Start(":8080"))
}
