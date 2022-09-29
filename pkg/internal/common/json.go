package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func GetJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	jsonData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		return err
	}
	return nil
}
