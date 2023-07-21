package persona

import (
	"net/http"

	"github.com/String-xyz/go-lib/v2/common"
)

/*

The account represents a verified individual and contains one or more inquiries.
The primary use of the account endpoints is to fetch previously submitted information for an individual.

*/

func (c *PersonaClient) CreateAccount(payload InquiryPayload) (*Inquiry, error) {
	inquiry := &Inquiry{}
	err := c.doRequest(http.MethodPost, "/v1/accounts", payload, inquiry)
	if err != nil {
		return nil, common.StringError(err, "failed to create an account")
	}
	return inquiry, nil
}

func (c *PersonaClient) GetAccountById(id string) (*AccountResponse, error) {
	account := &AccountResponse{}
	err := c.doRequest(http.MethodGet, "/v1/accounts/"+id, nil, account)
	if err != nil {
		return nil, common.StringError(err, "failed to get account")
	}
	return account, nil
}

func (c *PersonaClient) ListAccounts() (*ListAccountResponse, error) {
	accounts := &ListAccountResponse{}
	err := c.doRequest(http.MethodGet, "/v1/accounts", nil, accounts)
	if err != nil {
		return nil, common.StringError(err, "failed to list accounts")
	}
	return accounts, nil
}
