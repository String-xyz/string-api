package main

import (
	"log"
	"os"

	"github.com/String-xyz/string-api/src/api"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// TODO: create db connection to Postgres
	// Checkout https://github.com/jmoiron/sqlx
	dbuser := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := "user=" + dbuser + " dbname=" + dbname + " sslmode=disable"
	db, err := sqlx.Connect(
		"postgres", // requires a driver
		psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to db:", err)
	}
	// Query Builder https://github.com/Masterminds/squirrel

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
