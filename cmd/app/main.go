package main

import (
	"log"
	"os"

	"github.com/String-xyz/string-api/api"
	"github.com/joho/godotenv"
)

func main() {
	// TODO: create db connection to Postgres
	// Checkout https://github.com/jmoiron/sqlx
	// Query Builder https://github.com/Masterminds/squirrel

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// load db
	db := "test"

	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}

	// setup api
	err = api.Start(api.APIConfig{
		DB:   db,
		Port: port,
	})

}
