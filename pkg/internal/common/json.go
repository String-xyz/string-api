package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

func GetJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return StringError(err)
	}
	defer response.Body.Close()
	jsonData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return StringError(err)
	}
	targetType := reflect.TypeOf(target)
	if len(jsonData) != int(targetType.Size()) {
		return StringError(errors.New("Malformed JSON Response"))
	}
	err = json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		return StringError(err)
	}
	return nil
}
