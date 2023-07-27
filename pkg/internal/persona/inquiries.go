package persona

import (
	"net/http"

	"github.com/String-xyz/go-lib/v2/common"
)

/*
The inquiry represents a single instance of an individual attempting to verify their identity.
The primary use of the inquiry endpoints is to fetch submitted information from the flow.

Inquiries are created when the individual begins to verify their identity.
Check for the following statuses to determine whether the individual has finished the flow.

* Created	- The individual started the inquiry.
* Pending	- The individual submitted a verification within the inquiry.
* Completed	- The individual passed all required verifications within the inquiry.

Approved/Declined (Optional)
These are optional statuses applied by you to execute custom decisioning logic.
* Expired	- The individual did not complete the inquiry within 24 hours.
* Failed	- The individual exceeded the allowed number of verification attempts on the inquiry and cannot continue.
*/

type InquiryCreate struct {
	Attributes InquiryCreationAttributes `json:"attributes"`
}

type InquiryCreateRequest struct {
	Data InquiryCreate `json:"data"`
}

type Inquiry struct {
	Id            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    InquiryAttributes `json:"attributes"`
	Relationships Relationships     `json:"relationships"`
}

func (i Inquiry) GetType() string {
	return i.Type
}

type InquiryResponse struct {
	Data     Inquiry    `json:"data"`
	Included []Included `json:"included"`
}

type ListInquiryResponse struct {
	Data  []Inquiry `json:"data"`
	Links Link      `json:"links"`
}

func (c *PersonaClient) CreateInquiry(request InquiryCreateRequest) (*InquiryResponse, error) {
	inquiry := &InquiryResponse{}
	err := c.doRequest(http.MethodPost, "/v1/inquiries", request, inquiry)
	if err != nil {
		return nil, common.StringError(err, "failed to create inquiry")
	}
	return inquiry, nil
}

func (c *PersonaClient) GetInquiryById(id string) (*InquiryResponse, error) {
	inquiry := &InquiryResponse{}
	err := c.doRequest(http.MethodGet, "/v1/inquiries/"+id, nil, inquiry)
	if err != nil {
		return nil, common.StringError(err, "failed to get inquiry")
	}
	return inquiry, nil
}

func (c *PersonaClient) ListInquiriesByAccount(accountId string) (*ListInquiryResponse, error) {
	inquiries := &ListInquiryResponse{}
	err := c.doRequest(http.MethodGet, "/v1/inquiries?filter[account-id]="+accountId, nil, inquiries)
	if err != nil {
		return nil, common.StringError(err, "failed to list inquiries")
	}
	return inquiries, nil
}
