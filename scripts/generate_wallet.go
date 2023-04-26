package scripts

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"

	env "github.com/String-xyz/string-api/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// TODO: We could use the go=lib here
func StringError(err error, optionalMsg ...string) error {
	if err == nil {
		return nil
	}

	concat := ""

	for _, msgs := range optionalMsg {
		concat += msgs + " "
	}

	if errors.Cause(err) == nil || errors.Cause(err) == err {
		return errors.Wrap(errors.New(err.Error()), concat)
	}

	return errors.Wrap(err, concat)
}

// SSMPutParameterAPI defines the interface for the PutParameter function.
// We use this interface to test the function using a mocked service.
type SSMPutParameterAPI interface {
	PutParameter(ctx context.Context,
		params *ssm.PutParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

// SSMGetParameterAPI defines the interface for the GetParameter function.
// We use this interface to test the function using a mocked service.
type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

func AddStringParameter(c context.Context, api SSMPutParameterAPI, input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	res, err := api.PutParameter(c, input)
	if err != nil {
		return nil, StringError(err)
	}
	return res, nil
}

func FindParameter(c context.Context, api SSMGetParameterAPI, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	res, err := api.GetParameter(c, input)
	if err != nil {
		return nil, StringError(err)
	}
	return res, nil
}

func PutSSM(name string, value string, overwrite bool) error {
	if name == "" {
		return StringError(errors.New("unnamed parameter"))
	}
	if value == "" {
		return StringError(errors.New("empty value"))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return StringError(err)
	}
	ssmClient := ssm.NewFromConfig(cfg)
	keyId := env.Var.AWS_KMS_KEY_ID
	input := &ssm.PutParameterInput{
		Name:      &name,
		Value:     &value,
		Type:      types.ParameterTypeSecureString,
		Overwrite: &overwrite,
		KeyId:     &keyId,
	}
	_, err = AddStringParameter(context.TODO(), ssmClient, input)
	if err != nil {
		return StringError(err)
	}
	return nil
}

func GetSSM(name string) (string, error) {
	if name == "" {
		return "", StringError(errors.New("unnamed parameter"))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", StringError(err)
	}
	ssmClient := ssm.NewFromConfig(cfg)
	decrypt := true
	input := &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &decrypt,
	}
	results, err := FindParameter(context.TODO(), ssmClient, input)
	if err != nil {
		return "", StringError(err)
	}
	return *results.Parameter.Value, nil
}

func GenerateWallet() error {
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded

	preExistingWallet, _ := GetAddress()
	if preExistingWallet != "" {
		fmt.Printf("\n WARNING: WALLET CREDENTIALS FOR %+v ARE ALREADY BEING STORED IN SSM.  THIS SCRIPT WILL EXIT.", preExistingWallet)
		return nil
	}

	sk, err := crypto.GenerateKey()
	if err != nil {
		return StringError(err)
	}

	blob, err := EncryptStringToKMS(hexutil.Encode(crypto.FromECDSA(sk))[2:])
	if err != nil {
		return StringError(err)
	}
	err = PutSSM("string-encrypted-sk", string(blob), true)
	if err != nil {
		return StringError(err)
	}

	pk := sk.Public()
	pkECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return StringError(err)
	}
	addrStr := crypto.PubkeyToAddress(*pkECDSA).Hex()
	fmt.Printf("\nGenerated and encrypted new private key to SSM for wallet %+v", addrStr)

	err = PutSSM("string-wallet-address", addrStr, true)
	if err != nil {
		return StringError(err)
	}
	fmt.Printf("\nWallet address was also added to SSM.")
	return nil
}

func GetPrivateKey() (string, error) {
	blobStr, err := GetSSM("string-encrypted-sk")
	if err != nil {
		return "", StringError(err)
	}
	sk, err := DecryptBlobFromKMS(blobStr)
	if err != nil {
		return "", StringError(err)
	}
	return sk, nil
}

func GetAddress() (string, error) {
	address, err := GetSSM("string-wallet-address")
	if err != nil {
		return "", StringError(err)
	}
	return address, nil
}

func EncryptBytesToKMS(data []byte) (string, error) {
	region := env.Var.AWS_REGION
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", StringError(err)
	}
	kmsService := kms.New(session)
	keyId := env.Var.AWS_KMS_KEY_ID
	result, err := kmsService.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: data,
	})
	if err != nil {
		return "", StringError(err)
	}
	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

func EncryptStringToKMS(data string) (string, error) {
	res, err := EncryptBytesToKMS([]byte(data))
	if err != nil {
		return "", StringError(err)
	}
	return res, nil
}

func DecryptBlobFromKMS(blob string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return "", StringError(err)
	}
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", StringError(err)
	}
	kmsService := kms.New(session)
	result, err := kmsService.Decrypt(&kms.DecryptInput{CiphertextBlob: bytes})
	if err != nil {
		return "", StringError(err)
	}
	return string(result.Plaintext), nil
}
