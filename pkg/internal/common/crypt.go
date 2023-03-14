package common

import (
	"encoding/base64"
	"os"

	commonlib "github.com/String-xyz/go-lib/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func EncryptBytesToKMS(data []byte) (string, error) {
	region := os.Getenv("AWS_REGION")
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", commonlib.StringError(err)
	}
	kmsService := kms.New(session)
	keyId := os.Getenv("AWS_KMS_KEY_ID")
	result, err := kmsService.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: data,
	})
	if err != nil {
		return "", commonlib.StringError(err)
	}
	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

func EncryptStringToKMS(data string) (string, error) {
	res, err := EncryptBytesToKMS([]byte(data))
	if err != nil {
		return "", commonlib.StringError(err)
	}
	return res, nil
}

func DecryptBlobFromKMS(blob string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return "", commonlib.StringError(err)
	}
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", commonlib.StringError(err)
	}
	kmsService := kms.New(session)
	result, err := kmsService.Decrypt(&kms.DecryptInput{CiphertextBlob: bytes})
	if err != nil {
		return "", commonlib.StringError(err)
	}
	return string(result.Plaintext), nil
}
