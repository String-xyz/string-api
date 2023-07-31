package main

import (
	"log"
	"os"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	env "github.com/String-xyz/go-lib/v2/config"
	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// load env vars
	err := env.LoadEnv(&config.Var)
	if err != nil {
		panic(err)
	}
	lg := zerolog.New(os.Stdout)
	if !libcommon.IsLocalEnv() {
		setupTracer()
		defer profiler.Stop()
		defer tracer.Stop()
	}

	port := config.Var.PORT

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
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
		tracer.WithService("api"),
		tracer.WithEnv(config.Var.ENV),
	)

	err := profiler.Start(
		profiler.WithService("api"),
		profiler.WithEnv(config.Var.ENV),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		))

	if err != nil {
		log.Fatal(err)
	}
}
