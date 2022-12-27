package common

import (
	"context"
	"flag"

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

// AddStringParameter creates an AWS Systems Manager string parameter
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region
//	api is the interface that defines the method call
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a PutParameterOutput object containing the result of the service call and nil
//	Otherwise, nil and an error from the call to PutParameter
func AddStringParameter(c context.Context, api SSMPutParameterAPI, input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	res, err := api.PutParameter(c, input)
	if err != nil {
		return nil, StringError(err)
	}
	return res, nil
}

func PutSSM(name string, value string) error {
	paramName := flag.String("n", "", name)
	paramValue := flag.String("v", "", value)
	flag.Parse()
	if *paramName == "" {
		return StringError(errors.New("unnamed parameter"))
	}
	if *paramValue == "" {
		return StringError(errors.New("empty value"))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return StringError(err)
	}
	ssmClient := ssm.NewFromConfig(cfg)
	input := &ssm.PutParameterInput{
		Name:  paramName,
		Value: paramValue,
		Type:  types.ParameterTypeString,
	}
	_, err = AddStringParameter(context.TODO(), ssmClient, input)
	if err != nil {
		return StringError(err)
	}
	return nil
}
