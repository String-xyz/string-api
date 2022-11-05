package service

import (
	"math"
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/checkout/checkout-sdk-go"
	checkoutCommon "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
)

func getConfig() (*checkout.Config, error) {
	var sk = os.Getenv("CHECKOUT_SECRET_KEY")
	var pk = os.Getenv("CHECKOUT_PUBLIC_KEY")

	var config, err = checkout.SdkConfig(&sk, &pk, checkout.Sandbox)
	if err != nil {
		return nil, common.StringError(err)
	}
	return config, err
}

func convertAmount(amount float64) uint64 {
	return uint64(math.Round(amount * 100))
}

func CreateToken(card *tokens.Card) (*tokens.Response, error) {
	var config, err = getConfig()
	if err != nil {
		return nil, common.StringError(err)
	}
	client := tokens.NewClient(*config)

	res, err := client.Request(&tokens.Request{Card: card})
	if err != nil {
		return res, common.StringError(err)
	}
	return res, nil
}

type AuthorizedCharge struct {
	AuthID       string
	InstrumentID string
	Last4        string
}

func AuthorizeCharge(amount float64, userWallet string, tokenId string) (AuthorizedCharge, error) {
	res := AuthorizedCharge{}
	var config, err = getConfig()
	if err != nil {
		return res, common.StringError(err)
	}
	client := payments.NewClient(*config)

	// Generate a payment token ID in case we don't yet have one in the front end
	// For testing purposes only
	card := tokens.Card{
		Type:        checkoutCommon.Card,
		Number:      "4242424242424242",
		ExpiryMonth: 2,
		ExpiryYear:  2024,
		Name:        "Customer Name",
		CVV:         "100",
	}
	paymentToken, err := CreateToken(&card)
	if err != nil {
		return res, common.StringError(err)
	}
	paymentTokenID := paymentToken.Created.Token
	if tokenId != "" {
		paymentTokenID = tokenId
	}

	usd := convertAmount(amount)

	capture := false
	request := &payments.Request{
		Source: payments.TokenSource{
			Type:  checkoutCommon.Token.String(),
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
	response, err := client.Request(request, &params)

	if err != nil {
		return res, common.StringError(err)
	}

	// Collect authorization ID and Instrument ID
	res.AuthID = response.Processed.ID
	if response.Processed.Source.CardSourceResponse != nil {
		res.Last4 = response.Processed.Source.CardSourceResponse.Last4
		res.InstrumentID = response.Processed.Source.CardSourceResponse.ID
	}
	// TODO: Create entry for authorization in our DB associated with userWallet
	return res, nil
}

func CaptureCharge(amount float64, userWallet string, authorizationID string) (*payments.CapturesResponse, error) {
	var config, err = getConfig()
	if err != nil {
		return nil, common.StringError(err)
	}
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
	if err != nil {
		return nil, common.StringError(err)
	}

	// TODO: Create entry for capture in our DB associated with userWallet
	return res, nil
}
