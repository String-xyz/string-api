package main

import (
	"log"
	"os"

	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}
	db := store.GetPGInstance()
	// setup api
	_ = api.Start(api.APIConfig{
		DB:   db,
		Port: port,
	})
}
