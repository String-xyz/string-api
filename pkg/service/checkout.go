package service

import (
	"math"
	"os"

	"github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
)

func getConfig() (*checkout.Config, error) {
	var sk = os.Getenv("CHECKOUT_SECRET_KEY")
	var pk = os.Getenv("CHECKOUT_PUBLIC_KEY")

	var config, err = checkout.SdkConfig(&sk, &pk, checkout.Sandbox)
	if err != nil {
		return nil, err
	}
	return config, err
}

func convertAmount(amount float64) uint64 {
	return uint64(math.Round(amount * 100))
}

func Token(card *tokens.Card) (*tokens.Response, error) {
	var config, err = getConfig()
	client := tokens.NewClient(*config)

	res, err := client.Request(&tokens.Request{Card: card})
	return res, err
}

func Authorize(amount float64, userWallet string, tokenId string) (string, error) {
	var config, err = getConfig()
	client := payments.NewClient(*config)

	// Generate a payment token ID in case we don't yet have one in the front end
	// For testing purposes only
	card := tokens.Card{
		Type:        common.Card,
		Number:      "4242424242424242",
		ExpiryMonth: 2,
		ExpiryYear:  2024,
		Name:        "Customer Name",
		CVV:         "100",
	}
	paymentToken, err := Token(&card)
	if err != nil {
		return "", err
	}
	paymentTokenID := paymentToken.Created.Token
	if tokenId != "" {
		paymentTokenID = tokenId
	}

	usd := convertAmount(amount)

	capture := false
	request := &payments.Request{
		Source: payments.TokenSource{
			Type:  common.Token.String(),
			Token: paymentTokenID,
		},
		Amount:   usd,
		Currency: "USD",
		Customer: &payments.Customer{
			Name: userWallet,
		},
		AuthorizationType: "Estimated",
		Capture:           &capture,
	}

	idempotencyKey := checkout.NewIdempotencyKey()
	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}
	res, err := client.Request(request, &params)

	if err != nil {
		return "", err
	}
	// TODO: Create entry for authorization in our DB associated with userWallet
	return res.Processed.ID, err
}

func Capture(amount float64, userWallet string, authorizationID string) (*payments.CapturesResponse, error) {
	var config, err = getConfig()
	client := payments.NewClient(*config)

	usd := convertAmount(amount)

	idempotencyKey := checkout.NewIdempotencyKey()
	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}
	request := payments.CapturesRequest{
		Amount: usd,
	}
	res, err := client.Captures(authorizationID, &request, &params)
	// TODO: Create entry for capture in our DB associated with userWallet
	return res, err
}
