//go:build integration
// +build integration

package persona

import (
	"fmt"
	"testing"

	env "github.com/String-xyz/go-lib/v2/config"
	"github.com/stretchr/testify/assert"

	"github.com/String-xyz/string-api/config"
)

func init() {
	err := env.LoadEnv(&config.Var, "../../../.env")
	if err != nil {
		fmt.Printf("error loading env: %v", err)
	}
}

func client() *PersonaClient {
	return New(config.Var.PERSONA_API_KEY)
}

func TestIntegrationCreateAccount(t *testing.T) {
	request := AccountCreateRequest{AccountCreate{Attributes: CommonFields{NameFirst: "Mister", NameLast: "Tester"}}}
	account, err := client().CreateAccount(request)
	assert.NoError(t, err)
	assert.NotNil(t, account)
}

func TestIntegrationGetAccount(t *testing.T) {
	account, err := client().GetAccountById("act_Q1zEPYBZ6Qx8qJKcMrwDXxVA")
	assert.NoError(t, err)
	assert.NotNil(t, account)
}

func TestIntegrationListAccounts(t *testing.T) {
	accounts, err := client().ListAccounts()
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.NotEmpty(t, accounts.Data)
}

func TestIntegrationCreateInquiry(t *testing.T) {
	request := InquiryCreateRequest{InquiryCreate{Attributes: InquiryCreationAttributes{AccountId: "act_ndJNqdhWNi44S4Twf4bqzod1", InquityTemplateId: "itmpl_z2so7W2bCFHELp2dhxqqQjGy"}}}
	inquiry, err := client().CreateInquiry(request)
	assert.NoError(t, err)
	assert.NotNil(t, inquiry)
}

func TestIntegrationGetInquiry(t *testing.T) {
	inquiry, err := client().GetInquiryById("inq_kmXCg5pLzWTwg2LuAjiaBsoC")
	assert.NoError(t, err)
	assert.NotNil(t, inquiry)
}

func TestIntegrationListInquiriesByAccount(t *testing.T) {
	inquiries, err := client().ListInquiriesByAccount("act_ndJNqdhWNi44S4Twf4bqzod1")
	assert.NoError(t, err)
	assert.NotNil(t, inquiries)
	assert.NotEmpty(t, inquiries.Data)
}

func TestIntegrationGetVerification(t *testing.T) {
	verification, err := client().GetVerificationById("ver_ww2rkwtA6c9FiuCG8Jsk1DJt")
	assert.NoError(t, err)
	assert.NotNil(t, verification)
}
