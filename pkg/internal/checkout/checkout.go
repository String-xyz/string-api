package checkout

import (
	"errors"
	"math/rand"
	"time"

	"github.com/String-xyz/go-lib/v2/common"
	instruments "github.com/checkout/checkout-sdk-go/instruments/nas"
	"github.com/checkout/checkout-sdk-go/nas"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
	"github.com/rs/zerolog/log"

	"github.com/String-xyz/string-api/config"
)

type Checkout struct {
	ckoAPI   *nas.Api
	Payment  Payments
	Customer Customers
	Events   Events
}

func New() *Checkout {
	client := defaultAPI()
	return &Checkout{
		Payment:  Payments{client},
		Customer: Customers{client},
		Events:   Events{client},
	}
}

func (p Payments) Authorize(source Source, request PaymentRequest) (*PaymentResponse, error) {
	switch source.GetType() {
	case payments.TokenSource:
		token := source.(TokenSource)
		return p.AuthorizeWithToken(token, request)
	case payments.IdSource:
		id := source.(IdSource)
		return p.AuthorizeWithId(id, request)
	}
	return nil, common.StringError(errors.New("invalid source type, please provide either token or id type"))
}

func (p Payments) AuthorizeWithToken(token TokenSource, request PaymentRequest) (*PaymentResponse, error) {
	request.Source = tokenToSource(token)
	resp, err := p.authorizing(request)
	if err != nil {
		log.Err(err).Msg("internal checkout error while authorizing payment with token")
		return nil, err
	}
	return resp, nil
}

// AuthorizeWithId authorizes a payment with an id(instrument) type and returns a PaymentResponse and an error if any.
func (p Payments) AuthorizeWithId(id IdSource, request PaymentRequest) (*PaymentResponse, error) {
	request.Source = idToSource(id)
	resp, err := p.authorizing(request)
	if err != nil {
		log.Err(err).Msg("internal checkout error while authorizing payment with id")
		return nil, err
	}
	return resp, nil
}

// AuthorizeWithCard authorizes a payment with a card source type and returns a PaymentResponse
// and an error if any.
// PCI compliance is required to use this method.
func (p Payments) AuthorizeWithCard(card CardSource, request PaymentRequest) (*PaymentResponse, error) {
	request.Source = cardToSource(card)
	resp, err := p.authorizing(request)
	if err != nil {
		log.Err(err).Msg("internal checkout error while authorizing payment with card")
		return nil, err
	}
	return resp, nil
}

// AuthorizeWithCustomer authorizes a payment with a customer source type and returns a PaymentResponse and an error if any.
// A default card is required to use this method.
func (p Payments) AuthorizeWithCustomer(request PaymentRequest) (*PaymentResponse, error) {
	resp, err := p.authorizing(request)
	if err != nil {
		log.Err(err).Msg("internal checkout error while authorizing payment with customerId")
		return nil, common.StringError(err)
	}
	return resp, nil
}

// authorizing authorizes a payment with a payment request and returns a PaymentResponse and an error if any.
func (p Payments) authorizing(request PaymentRequest) (*PaymentResponse, error) {
	resp, err := p.client.Payments.RequestPayment(request, nil)
	if err != nil {
		return nil, common.StringError(err)
	}
	return resp, nil
}

// Capture captures a payment with a paymentId and returns a CaptureResponse and an error if any.
// is important to note that the capturing of a payment happens asychronously and the response does not indicate
// a successful capture. The response only indicates that the request was successfully sent to the Checkout API.
// To know if the capture was successful, we need to listen to the webhook event that is sent to the webhook endpoint
func (p Payments) Capture(paymentId string, request CaptureRequest) (*CaptureResponse, error) {
	resp, err := p.client.Payments.CapturePayment(paymentId, request, nil)
	if err != nil {
		return nil, common.StringError(err)
	}
	return resp, nil
}

func (p Payments) GetById(paymentId string) (*GetPaymentResponse, error) {
	resp, err := p.client.Payments.GetPaymentDetails(paymentId)
	if err != nil {
		log.Err(err).Msg("internal checkout error while getting payment by id")
		return nil, common.StringError(err)
	}
	return resp, nil
}

func (c Customers) Create(request CustomerRequest) (*CreateResponse, error) {
	resp, err := c.client.Customers.Create(request)
	if err != nil {
		log.Err(err).Msg("internal checkout error while creating customer")
		return nil, common.StringError(err)
	}
	return resp, nil
}

// GetById gets a customer by id and returns a GetCustomerResponse and an error if any.
// This method also returns all the instruments associated with the customer.
func (c Customers) GetById(customerId string) (*CustomerResponse, error) {
	resp, err := c.client.Customers.Get(customerId)
	if err != nil {
		log.Err(err).Msg("internal checkout error while getting customer by id")
		return nil, common.StringError(err)
	}
	return resp, nil
}

// ListInstruments is a convenience method that gets all the instruments associated with a customer.
// It returns InstrumentList and an error if any.
func (c Customers) ListInstruments(customerId string) ([]CardInstrument, error) {
	resp, err := c.client.Customers.Get(customerId)
	if err != nil {
		log.Err(err).Msg("internal checkout error while getting the customers instruments")
		return []CardInstrument{}, common.StringError(err)
	}

	return hydrateCardInstrument(resp.Instruments), nil
}

func hydrateCardInstrument(resp []instruments.GetInstrumentResponse) []CardInstrument {
	var instruments []CardInstrument
	for _, instrument := range resp {
		card := instrument.GetCardInstrumentResponse
		instruments = append(instruments, CardInstrument{
			Id:          card.Id,
			Last4:       card.Last4,
			ExpiryMonth: card.ExpiryMonth,
			ExpiryYear:  card.ExpiryYear,
			Scheme:      card.Scheme,
			Type:        string(card.Type),
			CardType:    string(card.CardType),
		})
	}
	return instruments
}

// DevCardToken returns a token for a test card
func DevCardToken() string {
	request := tokens.CardTokenRequest{
		Type:        tokens.Card,
		Number:      getTestCard(config.Var.CARD_FAIL_PROBABILITY),
		ExpiryMonth: 10,
		ExpiryYear:  2025,
		Name:        "DEV TOKEN",
		CVV:         "123",
	}

	response, err := defaultAPI().Tokens.RequestCardToken(request)
	if err != nil {
		log.Err(err).Msg("internal checkout error while getting dev card token")
		return ""
	}
	return response.Token
}

func getTestCard(failProbability float64) string {
	rand.Seed(time.Now().UnixNano())
	// Generate a random number between 0 and 1
	random := rand.Float64()
	if random < failProbability {
		// Choose a random fail card
		index := rand.Intn(len(failCards))
		return failCards[index]
	}
	// Choose a random success card
	index := rand.Intn(len(successCards))
	return successCards[index]
}
