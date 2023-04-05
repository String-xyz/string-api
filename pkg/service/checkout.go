// TODO: Make this service instantiable

package service

import (
	"math"
	"os"
	"strings"

	libcommon "github.com/String-xyz/go-lib/common"
	customer "github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/checkout/checkout-sdk-go"
	checkoutCommon "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
)

func getConfig() (*checkout.Config, error) {
	var sk = os.Getenv("CHECKOUT_SECRET_KEY")
	var pk = os.Getenv("CHECKOUT_PUBLIC_KEY")
	var env = os.Getenv("CHECKOUT_ENV")
	checkoutEnv := checkout.Sandbox

	if env == "prod" {
		checkoutEnv = checkout.Production
	}

	var config, err = checkout.SdkConfig(&sk, &pk, checkoutEnv)
	if err != nil {
		return nil, libcommon.StringError(err)
	}
	return config, err
}

func convertAmount(amount float64) uint64 {
	return uint64(math.Round(amount * 100))
}

func CreateToken(card *tokens.Card) (token *tokens.Response, err error) {
	config, err := getConfig()
	if err != nil {
		return nil, libcommon.StringError(err)
	}
	client := tokens.NewClient(*config)

	token, err = client.Request(&tokens.Request{Card: card})
	if err != nil {
		return token, libcommon.StringError(err)
	}
	return token, nil
}

func GetCustomerInstruments(Id string) ([]customer.CustomerInstrument, error) {
	config, err := getConfig()
	if err != nil {
		return nil, libcommon.StringError(err)
	}

	customer := customer.NewCustomer(*config)

	response, err := customer.GetCustomer(Id)
	if err != nil {
		return nil, libcommon.StringError(err)
	}

	if response.StatusResponse.StatusCode == 200 {
		return response.Customer.Instruments, nil
	}

	return nil, nil
}

type AuthorizedCharge struct {
	AuthId              string
	CheckoutFingerprint string
	Last4               string
	Issuer              string
	Approved            bool
	Status              string
	Summary             string
	CardType            string
	CardholderName      string
}

func AuthorizeCharge(p transactionProcessingData) (transactionProcessingData, error) {
	auth := AuthorizedCharge{}
	config, err := getConfig()
	if err != nil {
		return p, libcommon.StringError(err)
	}
	client := payments.NewClient(*config)

	var paymentTokenId string
	var paymentSource interface{}
	// if p.executionRequest.CardSourceId != "" {
	// 	paymentSource = payments.IDSource{
	// 		Type: "id",
	// 		ID:   p.executionRequest.CardSourceId,
	// 		CVV:  p.executionRequest.CVV,
	// 	}
	// } else {
	if p.executionRequest.CardToken != "" {
		paymentTokenId = p.executionRequest.CardToken
	} else if libcommon.IsLocalEnv() {

		// Generate a payment token ID in case we don't yet have one in the front end
		// For testing purposes only
		card := tokens.Card{
			Type:   checkoutCommon.Card,
			Number: "4242424242424242", // Success
			// Number: "4273149019799094", // succeed authorize, fail capture
			// Number: "4544249167673670", // Declined - Insufficient funds
			// Number:      "5148447461737269", // Invalid transaction (debit card)
			ExpiryMonth: 2,
			ExpiryYear:  2024,
			Name:        "Customer Name",
			CVV:         "100",
		}
		paymentToken, err := CreateToken(&card)
		if err != nil {
			return p, libcommon.StringError(err)
		}
		paymentTokenId = paymentToken.Created.Token
	}
	paymentSource = payments.TokenSource{
		Type:  checkoutCommon.Token.String(),
		Token: paymentTokenId,
	}
	// }

	fullName := p.user.FirstName + " " + p.user.MiddleName + " " + p.user.LastName
	fullName = strings.Replace(fullName, "  ", " ", 1) // If no middle name, ensure there is only one space between first name and last name

	usd := convertAmount(p.executionRequest.TotalUSD)
	capture := false
	request := &payments.Request{
		Source:   &paymentSource,
		Amount:   usd,
		Currency: "USD",
		Customer: &payments.Customer{
			Name:  fullName,
			Email: p.user.Email, // Replace with more robust email from platform and user
		},
		Capture:   &capture,
		PaymentIP: p.transactionModel.IPAddress,
	}

	idempotencyKey := checkout.NewIdempotencyKey()
	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}
	response, err := client.Request(request, &params)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	// Collect authorization ID and Instrument ID
	if response.Processed != nil {
		auth.AuthId = response.Processed.ID
		auth.Approved = *response.Processed.Approved
		auth.Status = string(response.Processed.Status)
		auth.Summary = response.Processed.ResponseSummary
		auth.CardType = string(response.Processed.Source.CardType)

		if response.Processed.Source.CardSourceResponse != nil {
			auth.Last4 = response.Processed.Source.CardSourceResponse.Last4
			auth.Issuer = response.Processed.Source.Issuer
			auth.CheckoutFingerprint = response.Processed.Source.CardSourceResponse.Fingerprint
			auth.CardholderName = response.Processed.Source.CardSourceResponse.Name
		}
	}
	p.cardAuthorization = &auth
	// TODO: Create entry for authorization in our DB associated with userWallet
	return p, nil
}

func CaptureCharge(p transactionProcessingData) (transactionProcessingData, error) {
	config, err := getConfig()
	if err != nil {
		return p, libcommon.StringError(err)
	}
	client := payments.NewClient(*config)

	usd := convertAmount(p.executionRequest.Quote.TotalUSD)

	idempotencyKey := checkout.NewIdempotencyKey()
	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}
	request := payments.CapturesRequest{
		Amount: usd,
	}

	capture, err := client.Captures(p.cardAuthorization.AuthId, &request, &params)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	p.cardCapture = capture

	// TODO: call action, err = client.Actions(capture.Accepted.ActionId) in another service to check on

	// TODO: Create entry for capture in our DB associated with userWallet
	return p, nil
}
