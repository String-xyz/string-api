package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func EncodeToBase64(object interface{}) (string, error) {
	buffer, err := json.Marshal(object)
	if err != nil {
		return "", StringError(err)
	}
	return base64.StdEncoding.EncodeToString(buffer), nil
}

func DecodeFromBase64[T any](from string) (T, error) {
	var result *T = new(T)
	buffer, err := base64.StdEncoding.DecodeString(from)
	if err != nil {
		fmt.Printf("\nFROM=%+v", from)
		return *result, StringError(err)
	}
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		return *result, StringError(err)
	}
	return *result, nil
}
