package common

import (
	"encoding/json"
	"fmt"
	"testing"

	b64 "encoding/base64"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSignAndValidateString(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	encodedMessage := "VGhhbmsgeW91IGZvciB1c2luZyBTdHJpbmchIEJ5IHNpZ25pbmcgdGhpcyBtZXNzYWdlIHlvdSBhcmU6CgoxKSBBdXRob3JpemluZyBTdHJpbmcgdG8gaW5pdGlhdGUgb2ZmLWNoYWluIHRyYW5zYWN0aW9ucyBvbiB5b3VyIGJlaGFsZiwgaW5jbHVkaW5nIHlvdXIgYmFuayBhY2NvdW50LCBjcmVkaXQgY2FyZCwgb3IgZGViaXQgY2FyZC4KCjIpIENvbmZpcm1pbmcgdGhhdCB0aGlzIHdhbGxldCBpcyBvd25lZCBieSB5b3UuCgpUaGlzIHJlcXVlc3Qgd2lsbCBub3QgdHJpZ2dlciBhbnkgYmxvY2tjaGFpbiB0cmFuc2FjdGlvbiBvciBjb3N0IGFueSBnYXMuCgpOb25jZTogejhQVk4wSzlNZW5hcGNOSnY0V2xvNFhkM1gxV2lCSVE5UE94b0hPc1ZuWFFjN0tCOEV2NTZzOTAvZU1OR25kWE03S2JtblZiMU9EZDdLc3VuekZEZW9SWGdwcTBaYTliek94VGJ0dVdzQnpnYnZsb3RjQ2V5NWx3VzRHMm5uTT0="

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
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

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
