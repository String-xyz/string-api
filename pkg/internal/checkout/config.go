package checkout

import (
	cko "github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/configuration"
	"github.com/checkout/checkout-sdk-go/nas"
	"github.com/rs/zerolog/log"

	"github.com/String-xyz/string-api/config"
)

func ckoEnv() configuration.Environment {
	if config.Var.CHECKOUT_ENV == "prod" {
		return configuration.Production()
	}

	return configuration.Sandbox()
}

func defaultAPI() *nas.Api {
	api, err := cko.
		Builder().
		StaticKeys().
		WithPublicKey(config.Var.CHECKOUT_PUBLIC_KEY).
		WithSecretKey(config.Var.CHECKOUT_SECRET_KEY).
		WithEnvironment(ckoEnv()). // or Environment.PRODUCTION
		Build()

	if err != nil {
		log.Err(err).Msg("error getting a default Checkout API Client")
		return nil
	}

	return api
}
