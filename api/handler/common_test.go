package handler

import (
	"testing"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeModelInput(t *testing.T) {
	m := model.User{Id: "user_0923840923840923840923"}
	err := SanitizeIdInput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "0923840923840923840923", m.Id)
}

func TestSanitizeModelOutput(t *testing.T) {
	m := model.User{Id: "0923840923840923840923"}
	err := SanitizeIdOutput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "user_0923840923840923840923", m.Id)
}

func TestSanitizeRelationalModelInput(t *testing.T) {
	m := model.ContactToPlatform{ContactId: "contact_123", PlatformId: "platform_456"}
	err := SanitizeIdInput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "123", m.ContactId)
	assert.Equal(t, "456", m.PlatformId)
}

func TestSanitizeRelationalModelOutput(t *testing.T) {
	m := model.ContactToPlatform{ContactId: "123", PlatformId: "456"}
	err := SanitizeIdOutput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "contact_123", m.ContactId)
	assert.Equal(t, "platform_456", m.PlatformId)
}
