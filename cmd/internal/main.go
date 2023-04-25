package main

import (
	"os"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/env"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// load .env file
	err := env.LoadEnv()
	if err != nil {
		panic(err)
	}

	if !libcommon.IsLocalEnv() {
		tracer.Start()
		defer tracer.Stop()
	}

	port, err := env.Get("PORT")
	if err != nil {
		panic(err)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	db := store.MustNewPG()
	lg := zerolog.New(os.Stdout)

	redis := store.NewRedis()

	// setup api
	api.StartInternal(api.APIConfig{
		DB:     db,
		Redis:  redis,
		Port:   port,
		Logger: &lg,
	})
}
