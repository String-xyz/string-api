package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// questions:
// do you want to allow upserts of entities? or only new entity creation?
// what is the format/type of the data in the user.tags jsonb? unit21 requires key:value pairs
// what format are phone numbers saved in the String db?

func CreateEntity(data StringData) (unit21Id string, err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")
	url := os.Getenv("UNIT21_URL") + "/entities/create"

	jsonBody := MapStringDataToEntity(data)

	reqBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Printf("Could not encode Entity to bytes: %s", err)
		return
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for Entity: %s", err)
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed to create Entity: %s", err)
		//handle 409 on update that is not allowed
		//handle 423, 500, 503 for retries
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting body from Entity create request: %s", err)
		return
	}

	bodystring := string(body)
	log.Printf("String of body from response: %s", bodystring)

	var entity *EntityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	log.Printf("Unit21Id: %s", entity.Unit21Id)
	return entity.Unit21Id, nil
}
