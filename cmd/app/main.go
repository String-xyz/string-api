package main

import (
	"log"
	"net/http"

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

	// e.GET("/", func(c echo.Context) error {
	// 	return transactHandler.Transact(c)
	// })

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
