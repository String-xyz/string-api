package unit21

// go run api/clients/unit21/apis/entities.go -data=value

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log" // do you have a logging pattern?
	"net/http"
	"os"
	"time"
)

type StringData struct { //temporary until the input data is established
	tags         map[string]string
	partnerName  string
	tier         string
	id           string
	userType     string
	status       string
	createdAt    int
	firstName    string
	middleName   string
	lastName     string
	emails       []string
	phones       []string
	ipAddresses  []string
	fingerprints []string
}

// questions:
// do you want to allow upserts? or only new entity creation?
// will an instrument exist first? and need to be attched in this call?
// what are whitelisted entities?
// will you want to run verification on an entity after creation?
// will there be batch uploads for entity creation?
func CreateEntity(data StringData) (unit21Id string, err error) {
	//where does tier come from? don't see it in the db
	apiKey := os.Getenv("UNIT21_API_KEY")
	url := os.Getenv("UNIT21_URL") + "/entities/create"

	//todo: move to a mapper file
	var userTagArr []string
	var userTags map[string]string //assumes the values are strings in the jsonb
	if data.tags != nil {
		for key, value := range data.tags {
			userTagArr = append(userTagArr, key+":"+value)
		}
	}
	userTagArr = append(userTagArr, "platform:"+data.partnerName) //partnerName is not in the schema
	userTagArr = append(userTagArr, "tier:"+data.tier)            // tier is not in the schema

	for key, value := range userTags {
		userTagArr = append(userTagArr, key+":"+value)
	}

	//todo: move above to mapper and call it in next line
	// https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go
	jsonBody := &NewEntity{
		GeneralData: &Entity{
			EntityId:      data.id,
			EntityType:    "user", //or employee?
			EntitySubType: data.userType,
			Status:        data.status,
			RegisteredAt:  data.createdAt, //probably need to convert to seconds
			Tags:          userTags,       // convert from jsonb into array of key:value string pairs
		},
		UserData: &User{
			FirstName:  data.firstName,
			MiddleName: data.middleName,
			LastName:   data.lastName,
		},
		CommunicationData: &Communication{
			Emails: data.emails, //might need to be converted to []string
			Phones: data.phones, //might need to be converted to []string
		},
		DigitalData: &DigitalInfo{
			IpAddresses:        data.ipAddresses,  //might need to be converted to []string
			ClientFingerprints: data.fingerprints, //schema doesn't have a fingerprint, might need to be convered to []string
		},
		InstrumentIds: data, //might need to be converted to []string
		CustomData: &Custom{
			Platform: data.partnerName, //partnerName is not in the schema
			Tier:     data.tier,        //tier is not in the schema
		},
		// add WorkflowOptions if not default?
	}
	bodyReader := bytes.NewReader(jsonBody)

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
	log.Printf("Response object: %s", res)

	// defer res.Body.Close() do we need this when doing encoding?
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

	log.Printf("Entity response: %s", entity)
	return entity.Unit21Id, nil
}
