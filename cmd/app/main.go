package main //why isn't this at the root of the project?

import (
	"log"
	"os"

	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file") // TODO: figure out why this wasnt hit
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}
	// if you are running local make sure to have an instance of pg running
	// this call will panic if it cant connect
	db := store.MustNewPG()
	lg := zerolog.New(os.Stdout)
	// setup api
	api.Start(api.APIConfig{
		DB:     db,
		Redis:  store.NewRedisStore(),
		Port:   port,
		Logger: &lg,
	})
}
