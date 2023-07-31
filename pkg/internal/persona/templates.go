package persona

import (
	"fmt"
	"net/http"
)

type Template struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

func (c *PersonaClient) GetTemplates() ([]Template, error) {
	var templates []Template
	err := c.doRequest(http.MethodGet, "/v1/templates", nil, &templates)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}
	return templates, nil
}

func (c *PersonaClient) GetTemplate(id string) (*Template, error) {
	template := &Template{}
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/templates/%s", id), nil, template)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return template, nil
}
