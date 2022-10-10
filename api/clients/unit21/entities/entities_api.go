package unit21

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log" // do you have a logging pattern yet that i should follow?
	"net/http"
	"os"
	"time"
)

// questions:
// do you want to allow upserts? or only new entity creation?
// will an instrument exist first? and need to be attched in this call?
// what are whitelisted entities?
// will you want to run verification on an entity after creation?
// will there be batch uploads for entity creation?
// where does tier come from? don't see it in the db
// where is partner name? it's not on the platform
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
