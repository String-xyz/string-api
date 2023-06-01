package checkout

import (
	"fmt"
	"testing"

	"github.com/String-xyz/string-api/config"
	"github.com/checkout/checkout-sdk-go/tokens"
)

func init() {
	err := config.LoadEnv("../../../.env")
	if err != nil {
		fmt.Printf("error loading env: %v", err)
	}
}

func generateToken() string {
	request := tokens.CardTokenRequest{
		Type:        tokens.Card,
		Number:      "4242424242424242",
		ExpiryMonth: 10,
		ExpiryYear:  2025,
		Name:        "Name",
		CVV:         "123",
	}

	response, err := defaultAPI().Tokens.RequestCardToken(request)
	if err != nil {
		fmt.Println(err)
	}
	return response.Token
}

func TestAuthorizeWithToken(t *testing.T) {
	p := &Payments{defaultAPI()}
	token := TokenSource{Token: generateToken()}
	resp, err := p.AuthorizeWithToken(token, PaymentRequest{Capture: true, Currency: "USD", Amount: 1000})
	if err != nil {
		t.Errorf("authorizeWithToken returned an error: %v", err)
	}
	if resp == nil {
		t.Errorf("authorizeWithToken returned a nil response")
	}
}

func TestAuthorizeWithCard(t *testing.T) {
	card := CardSource{
		Number:      "4242424242424242",
		ExpiryMonth: 10,
		ExpiryYear:  2025,
		Name:        "Name",
		Cvv:         "123",
	}
	p := &Payments{defaultAPI()}
	resp, err := p.AuthorizeWithCard(card, PaymentRequest{Capture: true, Currency: "USD", Amount: 1000})

	if err != nil {
		t.Errorf("authorizeWithCard returned an error: %v", err)
	}
	if resp == nil {
		t.Errorf("authorizeWithCard returned a nil response")
	}
}

func TestAuthorizeWithId(t *testing.T) {
	p := &Payments{defaultAPI()}
	tokenResp, err := p.AuthorizeWithToken(TokenSource{Token: generateToken(), StoreForFutureUse: true}, PaymentRequest{Capture: true, Currency: "USD", Amount: 1000})
	if err != nil {
		t.Errorf("authorizeWithToken returned an error: %v", err)
	}
	instrumentId := tokenResp.Source.ResponseCardSource.Id
	resp, err := p.AuthorizeWithId(IdSource{Id: instrumentId}, PaymentRequest{Capture: true, Currency: "USD", Amount: 1000})
	if err != nil {
		t.Errorf("authorizeWithId returned an error: %v", err)
	}
	if resp == nil {
		t.Errorf("authorizeWithId returned a nil response")
	}
}

func TestCapture(t *testing.T) {
	p := &Payments{defaultAPI()}
	tokenResp, err := p.AuthorizeWithToken(TokenSource{Token: generateToken()}, PaymentRequest{Capture: false, Currency: "USD", Amount: 1000})
	if err != nil {
		t.Errorf("authorizeWithToken returned an error: %v", err)
	}
	resp, err := p.Capture(tokenResp.Id, CaptureRequest{Amount: 1000})
	if err != nil {
		t.Errorf("capture returned an error: %v", err)
	}
	if resp == nil {
		t.Errorf("capture returned a nil response")
	}
}
