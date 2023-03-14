package common

import (
	"encoding/base64"
	"encoding/json"

	commonlib "github.com/String-xyz/go-lib/common"
)

func EncodeToBase64(object interface{}) (string, error) {
	buffer, err := json.Marshal(object)
	if err != nil {
		return "", commonlib.StringError(err)
	}
	return base64.StdEncoding.EncodeToString(buffer), nil
}

func DecodeFromBase64[T any](from string) (T, error) {
	var result *T = new(T)
	buffer, err := base64.StdEncoding.DecodeString(from)
	if err != nil {
		return *result, commonlib.StringError(err)
	}
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		return *result, commonlib.StringError(err)
	}
	return *result, nil
}
