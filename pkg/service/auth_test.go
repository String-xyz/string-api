package service

import (
	"testing"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/test/mocks"
	"github.com/String-xyz/string-api/pkg/test/stubs"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	a := NewAuth(stubs.AuthStrategyRepo{}, repository.NewUser(mocks.MockedDB()), repository.NewContact(mocks.MockedDB()))
	m := model.User{ID: "id"}
	token, err := a.GenerateJWT(m)
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Token)
}
func TestValidate(t *testing.T) {
	a := NewAuth(stubs.AuthStrategyRepo{}, repository.NewUser(mocks.MockedDB()), repository.NewContact(mocks.MockedDB()))
	m := model.User{ID: "id"}
	token, err := a.GenerateJWT(m)
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Token)
	valid, err := a.(AuthValidator).Validate(token.Token)
	assert.NoError(t, err)
	assert.True(t, valid)
}
