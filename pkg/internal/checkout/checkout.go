package checkout

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/checkout/checkout-sdk-go/nas"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/rs/zerolog/log"
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
func (c Customers) ListInstruments(customerId string) (InstrumentList, error) {
	resp, err := c.client.Customers.Get(customerId)
	if err != nil {
		log.Err(err).Msg("internal checkout error while getting the customers instruments")
		return resp.Instruments, common.StringError(err)
	}

	return resp.Instruments, nil
}

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
