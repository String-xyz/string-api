package common

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/pkg/errors"
)

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

func PutSSM(name string, value string) error {
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
	input := &ssm.PutParameterInput{
		Name:  &name,
		Value: &value,
		Type:  types.ParameterTypeString,
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
	input := &ssm.GetParameterInput{
		Name: &name,
	}
	results, err := FindParameter(context.TODO(), ssmClient, input)
	if err != nil {
		return "", StringError(err)
	}
	return *results.Parameter.Value, nil
}
