package persona

import (
	"fmt"
	"net/http"
)

type Inquiry struct {
	Id             string `json:"id"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	LastUpdatedAt  string `json:"last_updated_at"`
	CompletedSteps []struct {
		Type   string `json:"type"`
		Status string `json:"status"`
	} `json:"completed_steps"`
}

type InquiryPayload struct {
	AccountID string `json:"account_id"`
	Template  string `json:"template"`
}

func (c *PersonaClient) CreateInquiry(payload InquiryPayload) (*Inquiry, error) {
	inquiry := &Inquiry{}
	err := c.doRequest(http.MethodPost, "/v1/inquiries", payload, inquiry)
	if err != nil {
		return nil, fmt.Errorf("failed to create inquiry: %w", err)
	}
	return inquiry, nil
}

func (c *PersonaClient) GetInquiry(id string) (*Inquiry, error) {
	inquiry := &Inquiry{}
	err := c.doRequest(http.MethodGet, "/v1/inquiries/"+id, nil, inquiry)
	if err != nil {
		return nil, fmt.Errorf("failed to get inquiry: %w", err)
	}
	return inquiry, nil
}
