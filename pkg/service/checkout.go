// TODO: Make this service instantiable

package service

import (
	"fmt"
	"math"
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"

	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/internal/checkout"
	"github.com/String-xyz/string-api/pkg/model"
)

type AuthorizedCharge struct {
	PaymentId           string
	SourceId            string
	CheckoutFingerprint string
	Last4               string
	Issuer              string
	Approved            bool
	Status              string
	Summary             string
	CardType            string
	CardholderName      string
}

func convertAmount(amount float64) uint64 {
	return uint64(math.Round(amount * 100))
}

// Create Customer from the internal user and update
func createCustomer(user model.UserWithContact, platformId string) (string, error) {
	client := checkout.New()
	name := fmt.Sprintf("%s %s %s", user.FirstName, user.MiddleName, user.LastName)
	fullName := strings.Replace(name, "  ", " ", 1)

	resp, err := client.Customer.Create(checkout.CustomerRequest{
		Email: user.Email,
		Name:  fullName,
		Metadata: map[string]interface{}{
			"platformId": platformId,
			"internalId": user.Id,
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Error creating checkout customer from internal user")
		return "", err
	}

	return resp.Id, nil
}

func AuthorizeCharge(p transactionProcessingData) (transactionProcessingData, error) {
	usd := convertAmount(p.floatEstimate.TotalUSD)
	source := sourceForRequest(p)
	request := checkout.PaymentRequest{
		Amount:    int64(usd),
		Currency:  "USD",
		Capture:   false,
		PaymentIp: p.transactionModel.IPAddress,
	}
	request.Customer = customerForRequest(p)

	client := checkout.New()
	resp, err := client.Payment.Authorize(source, request)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	p.cardAuthorization, err = hydrateAuthorization(resp)

	return p, err
}

// CaptureCharge captures the payment for the given payment Id and sets the payment status to p
// If the payment is successful, the payment status will be "captured".
func CaptureCharge(p transactionProcessingData) (transactionProcessingData, error) {
	usd := convertAmount(p.floatEstimate.TotalUSD)
	client := checkout.New()
	captResp, err := client.Payment.Capture(p.cardAuthorization.PaymentId, checkout.CaptureRequest{Amount: int64(usd)})
	if err != nil {
		return p, libcommon.StringError(err)
	}

	if captResp.HttpMetadata.StatusCode != 202 {
		return p, libcommon.StringError(errors.Newf("capture failed with status code %d", captResp.HttpMetadata.StatusCode))
	}
	// Lets get the payment status to see if the payment was successful or not.
	// If the payment was successful, the payment status will be "captured".
	// If status is pending, we need to check the payment status again after a few seconds or use webhooks
	// which is not implemented yet.
	payResp, err := client.Payment.GetById(p.cardAuthorization.PaymentId)
	if err != nil {
		return p, libcommon.StringError(err)
	}

	if payResp.HttpMetadata.StatusCode != 200 {
		return p, libcommon.StringError(errors.Newf("get payment failed with status code %d", payResp.HttpMetadata.StatusCode))
	}

	p.PaymentStatus = payResp.Status
	p.ActionId = captResp.ActionId

	return p, nil
}

// sourceForRequest creates a sourcce from the transaction processing transactionProcessingData
// this source is needed for the checkout payment request, keep in mind that as of now
// only 2 sources are allowed; Id and Token, where Id is an instrument.
// Becase the return type is an interface, check for null when calling this function
// given that a nil value is return it none of the 2 data options is available.
func sourceForRequest(p transactionProcessingData) checkout.Source {
	paymentInfo := p.executionRequest.PaymentInfo
	if paymentInfo.CardId != nil && *paymentInfo.CardId != "" {
		return checkout.IdSource{
			BaseSource: checkout.BaseSource{Type: checkout.SourceTypeId},
			Id:         *paymentInfo.CardId,
			CVV:        *paymentInfo.CVV,
		}
	}

	if paymentInfo.CardToken != nil && *paymentInfo.CardToken != "" {
		return checkout.TokenSource{
			BaseSource:        checkout.BaseSource{Type: checkout.SourceTypeToken},
			Token:             *paymentInfo.CardToken,
			StoreForFutureUse: paymentInfo.SaveCard,
		}
	}

	// for local development, we can use a test card token
	if paymentInfo.CardToken != nil && *paymentInfo.CardToken == "" && config.Var.ENV == "local" {
		return checkout.TokenSource{
			BaseSource:        checkout.BaseSource{Type: checkout.SourceTypeToken},
			Token:             checkout.DevCardToken(),
			StoreForFutureUse: paymentInfo.SaveCard,
		}
	}

	return nil
}

func customerForRequest(p transactionProcessingData) *checkout.Customer {
	checkoutId := p.user.CheckoutId
	return &checkout.Customer{
		Id: checkoutId,
	}
}

func hydrateAuthorization(resp *checkout.PaymentResponse) (*AuthorizedCharge, error) {
	card := resp.Source.ResponseCardSource
	if card == nil {
		return nil, libcommon.StringError(errors.New("card source not found in response"))
	}

	return &AuthorizedCharge{
		PaymentId: resp.Id,
		Approved:  resp.Approved,
		SourceId:  card.Id,
		Issuer:    card.Issuer,
		Last4:     card.Last4,
	}, nil
}
