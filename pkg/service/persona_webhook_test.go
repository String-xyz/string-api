package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/String-xyz/string-api/pkg/test/data"
)

func TestHandle(t *testing.T) {
	webhook := personaWebhook{}

	tests := []struct {
		name string
		json string
		err  error
	}{
		{
			name: "Test Account event",
			json: data.PersonAccountJSON,
			err:  nil,
		},
		{
			name: "Test Inquiry event",
			json: data.PersonInquiryJSON,
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := webhook.Handle(context.Background(), []byte(tt.json))

			assert.Equal(t, tt.err, err)
		})
	}
}
