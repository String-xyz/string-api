package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func Encrypt(object interface{}, secret string) (string, error) {
	buffer, err := json.Marshal(object)
	if err != nil {
		return "", StringError(err)
	}
	return EncryptString(string(buffer), secret)
}

func Decrypt[T any](from string, secret string) (T, error) {
	var result T
	decrypted, err := DecryptString(from, secret)
	if err != nil {
		return result, StringError(err)
	}
	err = json.Unmarshal([]byte(decrypted), &result)
	if err != nil {
		return result, StringError(err)
	}
	return result, nil
}

func EncryptString(data string, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", StringError(err)
	}
	plainText := []byte(data)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", StringError(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptString(data string, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", StringError(err)
	}
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", StringError(err)
	}
	iv := cipherText[:aes.BlockSize]

	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func EncryptStringToKMS(data string) ([]byte, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return nil, StringError(err)
	}
	kmsService := kms.New(session)
	keyId := os.Getenv("AWS_KMS_KEY_ID")
	result, err := kmsService.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: []byte(data),
	})
	if err != nil {
		return nil, StringError(err)
	}
	return result.CiphertextBlob, nil
}

func DecryptBlobFromKMS(blob []byte) (string, error) {
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", StringError(err)
	}
	kmsService := kms.New(session)
	result, err := kmsService.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})
	if err != nil {
		return "", StringError(err)
	}
	return string(result.Plaintext), nil
}
