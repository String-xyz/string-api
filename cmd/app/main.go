package main

import (
	"os"

	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// load .env file
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	tracer.Start(
		tracer.WithServiceName("string-api"),
		tracer.WithEnv(os.Getenv("ENV")),
	)

	defer tracer.Stop()
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
