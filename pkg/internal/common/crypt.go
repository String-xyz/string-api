package common

import (
	"encoding/base64"

	"github.com/String-xyz/go-lib/common"
	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func EncryptBytesToKMS(data []byte) (string, error) {
	region, err := env.Get("AWS_REGION")
	if err != nil {
		return "", common.StringError(err)
	}
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", libcommon.StringError(err)
	}
	kmsService := kms.New(session)
	keyId, err := env.Get("AWS_KMS_KEY_ID")
	if err != nil {
		return "", common.StringError(err)
	}
	result, err := kmsService.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: data,
	})
	if err != nil {
		return "", libcommon.StringError(err)
	}
	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

func EncryptStringToKMS(data string) (string, error) {
	res, err := EncryptBytesToKMS([]byte(data))
	if err != nil {
		return "", libcommon.StringError(err)
	}
	return res, nil
}

func DecryptBlobFromKMS(blob string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return "", libcommon.StringError(err)
	}
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", libcommon.StringError(err)
	}
	kmsService := kms.New(session)
	result, err := kmsService.Decrypt(&kms.DecryptInput{CiphertextBlob: bytes})
	if err != nil {
		return "", libcommon.StringError(err)
	}
	return string(result.Plaintext), nil
}
