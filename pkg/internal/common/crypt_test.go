package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type randomObject1 struct {
	Timestamp int64  `json:"timestamp"`
	Email     string `json:"email"`
	Address   string `json:"address"`
}

func TestEncodeDecodeString(t *testing.T) {
	obj1 := "this is a test"

	obj1Encoded, err := EncodeToBase64(obj1)
	assert.NoError(t, err)

	obj1Decoded, err := DecodeFromBase64[string](obj1Encoded)
	assert.NoError(t, err)
	assert.Equal(t, obj1, obj1Decoded)
}

func TestEncodeDecodeObject(t *testing.T) {
	obj2 := randomObject1{Timestamp: time.Now().Unix(), Email: "test@test.com", Address: "0xdecafbabe"}

	obj2Encoded, err := EncodeToBase64(obj2)
	assert.NoError(t, err)

	obj2Decoded, err := DecodeFromBase64[randomObject1](obj2Encoded)
	assert.NoError(t, err)
	assert.Equal(t, obj2, obj2Decoded)
}

func TestEncryptDecryptString(t *testing.T) {
	str := "this is a string"

	strEncrypted, err := EncryptString(str, "secret_encryption_key_0123456789")
	assert.NoError(t, err)

	strDecrypted, err := DecryptString(strEncrypted, "secret_encryption_key_0123456789")
	assert.NoError(t, err)

	assert.Equal(t, str, strDecrypted)
}

func TestEncryptDecryptObject(t *testing.T) {
	obj := randomObject1{Timestamp: time.Now().Unix(), Email: "test@test.com", Address: "0xdecafbabe"}

	objEncoded, err := EncodeToBase64(obj)
	assert.NoError(t, err)

	objEncrypted, err := EncryptString(objEncoded, "secret_encryption_key_0123456789")
	assert.NoError(t, err)

	objDecrypted, err := DecryptString(objEncrypted, "secret_encryption_key_0123456789")
	assert.NoError(t, err)

	objDecoded, err := DecodeFromBase64[randomObject1](objDecrypted)
	assert.NoError(t, err)
	assert.Equal(t, obj, objDecoded)
}

func TestEncryptDecryptUnencoded(t *testing.T) {
	obj := randomObject1{Timestamp: time.Now().Unix(), Email: "test@test.com", Address: "0xdecafbabe"}

	objEncrypted, err := Encrypt(obj, "secret_encryption_key_0123456789")
	assert.NoError(t, err)

	objDecrypted, err := Decrypt[randomObject1](objEncrypted, "secret_encryption_key_0123456789")
	assert.NoError(t, err)
	assert.Equal(t, obj, objDecrypted)
}
