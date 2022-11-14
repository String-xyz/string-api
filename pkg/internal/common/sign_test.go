package common

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSignAndValidateString(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	obj1 := "0x44A4b9E2A69d86BA382a511f845CbF2E31286770"

	obj1Signed, err := EVMSign(obj1)
	assert.NoError(t, err)
	fmt.Printf("Signature: %+v", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, obj1)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
