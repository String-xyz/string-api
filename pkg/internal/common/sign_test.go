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
		Address:   "0xYourAddressHere",
		Timestamp: 123456789,
		Nonce:     "0xNonceGoesHere",
	}

	obj1Signed, err := EVMSign(obj1)
	assert.NoError(t, err)
	fmt.Printf("\nLogin Signature: %+v\n", obj1Signed)
	valid, err := ValidateExternalEVMSignature(obj1Signed, obj1.Address, obj1)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
