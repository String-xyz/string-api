package service

import (
	"testing"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	// Get address from SSM
	addr, err := GetAddress()
	assert.NoError(t, err)

	// Get private key from SSM
	sk, err := GetPrivateKey()
	assert.NoError(t, err)

	// Sign a message using private key
	msgToSign := []byte("The secrets to the universe")
	sig, err := common.EVMSignWithPrivateKey(msgToSign, sk, true)
	assert.NoError(t, err)

	// Verify the signature
	valid, err := common.ValidateExternalEVMSignature(sig, addr, msgToSign, true)
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
