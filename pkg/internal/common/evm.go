package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/lmittmann/w3"
)

func ParseParams(function *w3.Func, signature string, params []string) ([]byte, error) {
	signatureArgs := strings.Split(strings.Split(strings.Split(signature, "(")[1], ")")[0], ",")
	if len(signatureArgs) != len(params) {
		return nil, errors.New("executor parseParams: mismatched arguments")
	}
	args := []interface{}{}
	for i, s := range signatureArgs {
		switch s {
		case "address":
			args = append(args, w3.A(params[i]))
		case "bool":
			args = append(args, params[i] == "true" || params[i] == "TRUE")
		case "string":
			args = append(args, params[i])
		case "bytes":
			args = append(args, w3.B(params[i]))
		case "uint8":
			v, err := strconv.ParseUint(params[i], 0, 8)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		case "uint32":
			v, err := strconv.ParseUint(params[i], 0, 32)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		case "uint256":
			args = append(args, w3.I(params[i]))
		case "int8":
			v, err := strconv.ParseInt(params[i], 0, 8)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		case "int32":
			v, err := strconv.ParseInt(params[i], 0, 32)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		case "int256":
			args = append(args, w3.I(params[i]))
		default:
			return nil, errors.New("executor: parseParams: unsupported type")
		}
	}
	result, err := function.EncodeArgs(args)
	if err != nil {
		return nil, err
	}
	return result, nil
}
