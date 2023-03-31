package checkout

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/checkout/checkout-sdk-go"
	checkoutCommon "github.com/checkout/checkout-sdk-go/common"
	"github.com/checkout/checkout-sdk-go/customers"
	"github.com/checkout/checkout-sdk-go/httpclient"
	"github.com/checkout/checkout-sdk-go/instruments"
	"github.com/checkout/checkout-sdk-go/payments"
)

type Customer struct {
	API checkout.HTTPClient
}

// NewCustomer ...
func NewCustomer(config checkout.Config) *Customer {
	return &Customer{
		API: httpclient.NewClient(config),
	}
}

type CustomerResponse struct {
	StatusResponse *checkout.StatusResponse `json:"api_response,omitempty"`
	Customer       *CustomerData            `json:"customer,omitempty"`
}
type CustomerData struct {
	Id string `json:"id"`
	*customers.Customer
	Phone       *checkoutCommon.Phone `json:"phone,omitempty"`
	Metadata    map[string]string     `json:"metadata,omitempty"`
	Instruments []CustomerInstrument  `json:"instruments,omitempty"`
}

type CustomerInstrument struct {
	*payments.DestinationResponse
	*instruments.AccountHolder
}

func (c Customer) GetCustomer(customerId string) (*CustomerResponse, error) {
	url := fmt.Sprintf("/%v/%v", "customers", customerId)
	fmt.Printf("\n\n>>>>>>>>>>>> url: %+v", url)
	response, err := c.API.Get(url)
	resp := &CustomerResponse{
		StatusResponse: response,
	}
	if err != nil && response.StatusCode != http.StatusNotFound {
		fmt.Printf("\n\n>>>>>>>>>>>> err in GetCustomer: %+v", err)
		return resp, err
	}
	if response.StatusCode == http.StatusOK {
		var customer CustomerData
		err = json.Unmarshal(response.ResponseBody, &customer)
		resp.Customer = &customer
	}
	return resp, nil
}

// EXAMPLE DATA:
//
//	{
//		"id": "cus_y3oqhf46pyzuxjbcn2giaqnb44",
//		"email": "john.smith@example.com",
//		"default": "src_imu3wifxfvlebpqqq5usjrze6y",
//		"name": "John Smith",
//		"phone": {
//		  "country_code": "+1",
//		  "number": "5551234567"
//		},
//		"metadata": {
//		  "coupon_code": "NY2018",
//		  "partner_id": 123989
//		},
//		"instruments": [
//		  {
//			"id": "src_lmyvsjadlxxu7kqlgevt6ebkra",
//			"type": "card",
//			"fingerprint": "vnsdrvikkvre3dtrjjvlm5du4q",
//			"expiry_month": 6,
//			"expiry_year": 2025,
//			"name": "John Smith",
//			"scheme": "VISA",
//			"last4": "9996",
//			"bin": "454347",
//			"card_type": "Credit",
//			"card_category": "Consumer",
//			"issuer": "Test Bank",
//			"issuer_country": "US",
//			"product_id": "F",
//			"product_type": "CLASSIC",
//			"account_holder": {
//			  "billing_address": {
//				"address_line1": "123 Anywhere St.",
//				"address_line2": "Apt. 456",
//				"city": "Anytown",
//				"state": "AL",
//				"zip": "123456",
//				"country": "US"
//			  },
//			  "phone": {
//				"country_code": "+1",
//				"number": "5551234567"
//			  }
//			}
//		  }
//		]
//	  }
