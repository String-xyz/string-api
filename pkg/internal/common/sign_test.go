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

	obj1 := "Your Public Key Here"

	obj1Signed, err := EVMSign(obj1)
	assert.NoError(t, err)
	fmt.Printf("Signature: %+v", obj1Signed)
	valid, err := ValidateEVMSignature(obj1Signed, obj1)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
