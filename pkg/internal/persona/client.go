package persona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PersonaClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func New(apiKey string) *PersonaClient {
	return &PersonaClient{
		APIKey:  apiKey,
		BaseURL: "https://withpersona.com/api/",
		Client:  &http.Client{},
	}
}

func NewPersonaClient(baseURL, apiKey string) *PersonaClient {
	return &PersonaClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  &http.Client{},
	}
}

func (c *PersonaClient) doRequest(method, url string, payload, result interface{}) error {
	var body io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(method, c.BaseURL+url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Key-Inflection", "camel")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
