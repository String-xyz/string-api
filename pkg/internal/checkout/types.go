package checkout

import (
	ckocommon "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/customers"
	instruments "github.com/checkout/checkout-sdk-go/instruments/nas"
	ckonas "github.com/checkout/checkout-sdk-go/nas"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/payments/nas"
	"github.com/checkout/checkout-sdk-go/payments/nas/sources"
)

type (
	Payments struct {
		client *ckonas.Api
	}
	Instruments struct {
		client *ckonas.Api
	}
	Customers struct {
		client *ckonas.Api
	}
	Events struct {
		client *ckonas.Api
	}
)

const (
	SourceTypeId    = payments.IdSource
	SourceTypeToken = payments.TokenSource
)

type (
	SourceType         = payments.SourceType
	Customer           = ckocommon.CustomerRequest
	CustomerRequest    = customers.CustomerRequest
	CreateResponse     = ckocommon.IdResponse
	InstrumentList     = []instruments.GetInstrumentResponse
	CustomerResponse   = customers.GetCustomerResponse
	PaymentRequest     = nas.PaymentRequest
	PaymentResponse    = nas.PaymentResponse
	GetPaymentResponse = nas.GetPaymentResponse
	CaptureRequest     = nas.CaptureRequest
	CaptureResponse    = payments.CaptureResponse
	PaymentStatus      = payments.PaymentStatus
)

type CardInstrument struct {
	Id          string `json:"id,omitempty"`
	Cvv         string `json:"cvv,omitempty"`
	Last4       string `json:"last4,omitempty"`
	ExpiryMonth int    `json:"expiryMonth,omitempty"`
	ExpiryYear  int    `json:"expiryYear,omitempty"`
	Type        string `json:"type,omitempty"`
	CardType    string `json:"cardType,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
	Expired     bool   `json:"expired,omitempty"`
}

type Source interface {
	GetType() payments.SourceType
}

type BaseSource struct {
	Type payments.SourceType
}

type CardSource struct {
	BaseSource
	Number      string
	ExpiryMonth int
	ExpiryYear  int
	Name        string
	Cvv         string
	Stored      bool
}

type TokenSource struct {
	BaseSource
	// The token value returned by the Checkout.com API
	Token string
	// Wether the token can be used for future payments
	// by default this is set to true, but can be set to false
	// and if set to false, an instrument(which on this context is called source) will not be included in the response
	StoreForFutureUse bool
}

type IdSource struct {
	BaseSource
	// the id of the instrument, required.
	Id string
	// the cvv of the instrument if is a card, optional.
	CVV string
	// the payment method, only required for ACH, optional.
	PaymentMethod string
}

func cardToSource(card CardSource) Source {
	src := sources.NewRequestCardSource()
	src.Name = card.Name
	src.Number = card.Number
	src.Cvv = card.Cvv
	src.ExpiryMonth = card.ExpiryMonth
	src.ExpiryYear = card.ExpiryYear
	return src
}

func tokenToSource(token TokenSource) Source {
	src := sources.NewRequestTokenSource()
	src.Token = token.Token
	src.StoreForFutureUse = token.StoreForFutureUse
	return src
}

func idToSource(id IdSource) Source {
	src := sources.NewRequestIdSource()
	src.Id = id.Id
	src.Cvv = id.CVV
	return src
}

func (b BaseSource) GetType() payments.SourceType {
	return b.Type
}

// Test Credit Cards Only for Sandbox
var successCards = []string{
	"4242424242424242",
	"5436031030606378",
	"5305484748800098",
	"345678901234564",
}

var failCards = []string{
	"4644968546281686",
	"5355228287185489",
	"4546381219393284",
	"5355229757805879",
}
