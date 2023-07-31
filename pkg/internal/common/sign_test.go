package common

import (
	"encoding/json"
	"fmt"
	"testing"

	b64 "encoding/base64"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSignAndValidateString(t *testing.T) {

	encodedMessage := "Your base64 encoded String Here"

	// decode
	decoded, err := b64.URLEncoding.DecodeString(encodedMessage)
	assert.NoError(t, err)

	// sign
	obj1Signed, err := EVMSign(decoded, true)
	assert.NoError(t, err)
	fmt.Printf("\nString Signature: %+v\n", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, decoded, true)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}

func TestSignAndValidateStruct(t *testing.T) {
	// Paste the JSON output properties from whatever struct here
	obj1 := model.WalletSignaturePayload{
		Address:   "0xPasteYourAddressHere",
		Timestamp: 1010101010,
	}
	bytes, err := json.Marshal(obj1)
	assert.NoError(t, err)

	obj1Signed, err := EVMSign(bytes, true)
	assert.NoError(t, err)
	fmt.Printf("\nStruct Signature: %+v\n", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, bytes, true)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
