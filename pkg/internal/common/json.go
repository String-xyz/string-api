package common

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"time"

	commonlib "github.com/String-xyz/go-lib/common"
	"github.com/pkg/errors"
)

// For target interfaces which have a fixed size (no maps, arrays, etc)
func GetJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return commonlib.StringError(err)
	}
	defer response.Body.Close()
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return commonlib.StringError(err)
	}
	targetType := reflect.TypeOf(target)
	if len(jsonData) != int(targetType.Size()) {
		return commonlib.StringError(errors.New("Malformed JSON Response"))
	}
	err = json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		return commonlib.StringError(err)
	}
	return nil
}

// For target interfaces which do not have a fixed size (ie map[string]interface{} or a struct with an array in it)
func GetJsonGeneric(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return commonlib.StringError(err)
	}
	defer response.Body.Close()
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return commonlib.StringError(err)
	}
	err = json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		return commonlib.StringError(err)
	}
	return nil
}

func parseJSON[T any](b []byte) (T, error) {
	var r T
	if err := json.Unmarshal(b, &r); err != nil {
		return r, err
	}
	return r, nil
}
