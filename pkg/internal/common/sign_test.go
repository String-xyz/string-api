package common

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func SignAndValidateString(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	obj1 := "test string"

	obj1Signed, err := EVMSign(obj1)
	assert.NoError(t, err)
	valid, err := ValidateEVMSignature(obj1Signed, obj1)
	assert.NoError(t, err)
	assert.Equal(t, false, valid)
}
