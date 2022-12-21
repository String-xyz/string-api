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

	obj1 := "Your String Here"

	obj1Signed, err := EVMSign(obj1, true)
	assert.NoError(t, err)
	fmt.Printf("\nString Signature: %+v\n", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, obj1, true)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}

func TestSignAndValidateStruct(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	// Paste the JSON output properties from whatever struct here
	obj1 := model.WalletSignaturePayload{
		Address:   "0xPasteYourAddressHere",
		Timestamp: 1010101010,
	}

	obj1Signed, err := EVMSign(obj1, true)
	assert.NoError(t, err)
	fmt.Printf("\nStruct Signature: %+v\n", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, obj1, true)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
