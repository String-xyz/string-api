package unit21

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
)

func u21Put(url string, jsonBody any) (body []byte, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Printf("Could not encode %+v to bytes: %s", jsonBody, err)
		return nil, common.StringError(err)
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for %s: %s", url, err)
		return nil, common.StringError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed to update %s: %s", url, err)
		return nil, common.StringError(err)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting body from %s update request: %s", url, err)
		return nil, common.StringError(err)
	}

	log.Printf("String of body from response: %s", string(body))

	if res.StatusCode != 200 {
		log.Printf("Request failed to update %s: %s", url, fmt.Sprint(res.StatusCode))
		err = common.StringError(fmt.Errorf("request failed with status code %s and return body: %s", fmt.Sprint(res.StatusCode), string(body)))
		return
	}

	return body, nil
}

func u21Post(url string, jsonBody any) (body []byte, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Printf("Could not encode %+v to bytes: %s", jsonBody, err)
		return nil, common.StringError(err)
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for %s: %s", url, err)
		return nil, common.StringError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed to update %s: %s", url, err)
		return nil, common.StringError(err)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting body from %s update request: %s", url, err)
		return nil, common.StringError(err)
	}

	log.Printf("String of body from response: %s", string(body))

	if res.StatusCode != 200 {
		log.Printf("Request failed to update %s: %s", url, fmt.Sprint(res.StatusCode))
		err = common.StringError(fmt.Errorf("request failed with status code %s and return body: %s", fmt.Sprint(res.StatusCode), string(body)))
		return
	}

	return body, nil
}
