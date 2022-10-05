package service

import (
	"testing"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestQuote(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	tr := NewTransaction(nil)
	m := model.TransactionRequest{
		UserAddress: "0x44A4b9E2A69d86BA382a511f845CbF2E31286770",
		ChainID:     43113,
		CxAddr:      "0x861aF9Ed4fEe884e5c49E9CE444359fe3631418B",
		CxFunc:      "mintTo(address)",
		CxReturn:    "uint256",
		CxParams:    []string{"0x44A4b9E2A69d86BA382a511f845CbF2E31286770"},
		TxValue:     "0.08 eth",
		TxGasLimit:  "800000",
	}
	res, err := tr.Quote(m)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
func TestTransact(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	tr := NewTransaction(nil)
	m := model.TransactionRequest{
		UserAddress: "0x44A4b9E2A69d86BA382a511f845CbF2E31286770",
		ChainID:     43113,
		CxAddr:      "0x861aF9Ed4fEe884e5c49E9CE444359fe3631418B",
		CxFunc:      "mintTo(address)",
		CxReturn:    "uint256",
		CxParams:    []string{"0x44A4b9E2A69d86BA382a511f845CbF2E31286770"},
		TxValue:     "0.08 eth",
		TxGasLimit:  "800000",
	}
	req, err := tr.Quote(m)
	assert.NoError(t, err)

	txID, err := tr.Execute(req)
	assert.NoError(t, err)

	assert.NotEmpty(t, txID)
}
