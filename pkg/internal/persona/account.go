package persona

import (
	"net/http"

	"github.com/String-xyz/go-lib/v2/common"
)

/*
The account represents a verified individual and contains one or more inquiries.
The primary use of the account endpoints is to fetch previously submitted information for an individual.
*/

type Account struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes AccountAttributes `json:"attributes"`
}

type AccountCreate struct {
	Attributes CommonFields `json:"attributes"`
}

type AccountCreateRequest struct {
	Data AccountCreate `json:"data"`
}

type AccountResponse struct {
	Data Account `json:"data"`
}

type ListAccountResponse struct {
	Data  []Account `json:"data"`
	Links Link      `json:"links"`
}

func (c *PersonaClient) CreateAccount(request AccountCreateRequest) (*AccountResponse, error) {
	account := &AccountResponse{}
	err := c.doRequest(http.MethodPost, "/v1/accounts", request, account)
	if err != nil {
		return nil, common.StringError(err, "failed to create account")
	}
	return account, nil
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
