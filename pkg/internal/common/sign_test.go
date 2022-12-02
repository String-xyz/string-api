package common

import (
	"fmt"
	"testing"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSignAndValidateString(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	obj1 := "Your Public Key Here"

	obj1Signed, err := EVMSign(obj1)
	assert.NoError(t, err)
	fmt.Printf("\nPublic Key Signature: %+v\n", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, obj1)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}

func TestSignAndValidateStruct(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	// Paste the JSON output properties from wallet login request here
	obj1 := model.WalletSignaturePayload{
		Address:   "0x44A4b9E2A69d86BA382a511f845CbF2E31286770",
		Timestamp: 1670010656,
		Nonce:     "0x0b970ce1867862ba74daa65d83d16fb99e5cb88d0398e16f91593db1e4584b3e2e96cb65fab2c9ca070929d0308508012c2abe5ed3cdb5359327b4b19b15c4b701",
	}

	obj1Signed, err := EVMSign(obj1)
	assert.NoError(t, err)
	fmt.Printf("\nLogin Signature: %+v\n", obj1Signed)
	valid, err := ValidateExternalEVMSignature(obj1Signed, obj1.Address, obj1)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
