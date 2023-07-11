package persona

import (
	"fmt"
	"net/http"
)

type IdentityVerification struct {
	Id             string `json:"id"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	LastUpdatedAt  string `json:"last_updated_at"`
	CompletedSteps []struct {
		Type   string `json:"type"`
		Status string `json:"status"`
	} `json:"completed_steps"`
}

type IdentityVerificationPayload struct {
	AccountId string `json:"account_id"`
	Template  string `json:"template"`
}

type Verification struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func (c *PersonaClient) VerifyIdentity(payload IdentityVerificationPayload) (*IdentityVerification, error) {
	verification := &IdentityVerification{}
	err := c.doRequest(http.MethodPost, "/v1/verifications", payload, verification)
	if err != nil {
		return nil, fmt.Errorf("failed to verify identity: %w", err)
	}
	return verification, nil
}

func (c *PersonaClient) GetVerifications() ([]Verification, error) {
	var verifications []Verification
	err := c.doRequest(http.MethodGet, "/v1/verifications", nil, &verifications)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return verifications, nil
}
