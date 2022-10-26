package unit21

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
)

// type Unit21Data interface {
// 	Create(any) (unit21Id string, err error)
// 	Update(id string, updates any) (err error)
// }

// type unit21Data struct {
// 	repository repository.Transactable
// }

// func newUnit21Data(repo repository.Transactable) Unit21Data {
// 	return &unit21Data{repository: repo}
// }

// u21type == entities
// jsonBody = MapUserToEntity(user)
func create(datatype string, jsonBody any) (body []byte, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")
	url := os.Getenv("UNIT21_URL") + datatype + "/create"

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Printf("Could not encode %s to bytes: %s", datatype, err)
		return nil, common.StringError(err)
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for %s: %s", datatype, err)
		return nil, common.StringError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed to create %s: %s", datatype, err)
		//handle 409 on update that is not allowed
		//handle 423, 500, 503 for retries
		return nil, common.StringError(err)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting body from %s create request: %s", datatype, err)
		return nil, common.StringError(err)
	}

	log.Printf("String of body from response: %s", string(body))

	return body, nil
}

func update(datatype string, id string, jsonBody any) (body []byte, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")
	url := os.Getenv("UNIT21_URL") + "/" + datatype + "/" + id + "/update"

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Printf("Could not encode %s to bytes: %s", datatype, err)
		return nil, common.StringError(err)
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for %s: %s", datatype, err)
		return nil, common.StringError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed to update %s: %s", datatype, err)
		//handle 409 on update that is not allowed
		//handle 423, 500, 503 for retries
		return nil, common.StringError(err)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting body from %s update request: %s", datatype, err)
		return nil, common.StringError(err)
	}

	log.Printf("String of body from response: %s", string(body))

	return body, nil
}
