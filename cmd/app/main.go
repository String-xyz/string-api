package main

import (
	"log"
	"os"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// load .env file
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	lg := zerolog.New(os.Stdout)
	if !libcommon.IsLocalEnv() {
		setupTracer()
		defer profiler.Stop()
		defer tracer.Stop()
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// zerolog.SetGlobalLevel(zerolog.Disabled) // quiet mode
	db := store.MustNewPG()

	redis := store.NewRedis()

	// setup api
	api.Start(api.APIConfig{
		DB:     db,
		Redis:  redis,
		Port:   port,
		Logger: &lg,
	})
}

func setupTracer() {
	rules := []tracer.SamplingRule{tracer.RateRule(1)}
	tracer.Start(
		tracer.WithSamplingRules(rules),
		tracer.WithService("string-api"),
		tracer.WithEnv(os.Getenv("ENV")),
	)

	err := profiler.Start(
		profiler.WithService("string-api"),
		profiler.WithEnv(os.Getenv("ENV")),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		))

	if err != nil {
		log.Fatal(err)
	}
}
