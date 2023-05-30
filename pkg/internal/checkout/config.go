package checkout

import (
	cko "github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/configuration"
	"github.com/checkout/checkout-sdk-go/nas"
	"github.com/rs/zerolog/log"
)

const tst = "sk_sbox_qng7xllxv3oqj2dnhjs5tcf4zaa"

func ckoEnv() configuration.Environment {
	// if config.Var.CHECKOUT_ENV == "prod" {
	// 	return configuration.Production()
	// }

	return configuration.Sandbox()
}

func defaultAPI() *nas.Api {
	api, err := cko.
		Builder().
		StaticKeys().
		WithPublicKey("pk_sbox_wpmezqltqm4lc5jqmu6p7ccq6iu").
		WithSecretKey(tst).
		WithEnvironment(ckoEnv()). // or Environment.PRODUCTION
		Build()

	if err != nil {
		log.Err(err).Msg("error getting a default Checkout API Client")
		return nil
	}

	return api
}
