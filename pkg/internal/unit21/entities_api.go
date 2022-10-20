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
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Entity interface {
	Create(user model.User) (unit21Id string, err error)
	Update(id string, updates any) (err error)
	AddInstruments(entityId string, instrumentId []string) (err error)
}

type entity struct {
	userRepo repository.User
	// deviceRepo     repository.Device
	contactRepo repository.UserContact
	// instrumentRepo repository.Instrument
}

// With Device and Instrument
// func newEntity(user repository.User, device repository.Device, contact repository.UserContact, instrument repository.Instrument) Entity {
// 	return &entity{userRepo: user, deviceRepo: device, contactRepo: contact, instrumentRepo: instrument}
// }

func newEntity(user repository.User, contact repository.UserContact) Entity {
	return &entity{userRepo: user, contactRepo: contact}
}

// //
// u21entity := unit21.NewEntity()
// u21insturment := unit21.NewIntrument()

// u21entity.Create(user)
// //

// https://docs.unit21.ai/reference/create_entity
func (e entity) Create(user model.User) (unit21Id string, err error) {

	// ultimately may want a join here.
	// devices, err  := e.deviceRepo.ListUserID(user.ID, 100, 0)
	// instruments, err  := e.instrumentRepo.ListUserID(user.ID, 100, 0)

	// Get user contacts
	contacts, err := e.contactRepo.ListUserID(user.ID, 100, 0)
	if err != nil {
		log.Printf("Failed go get user contacts: %s", err)
		return "", common.StringError(err)
	}

	// Convert contact structs to an entity communication struct
	var communication communication
	for _, contact := range contacts {
		if contact.Type == "Email" {
			communication.Emails = append(communication.Emails, contact.Data)
		} else if contact.Type == "Phone" {
			communication.Phones = append(communication.Phones, contact.Data)
		}
	}
	log.Printf("communication: %s", communication)

	body, err := create("entities", MapUserToEntity(user, communication))
	if err != nil {
		log.Printf("Unit21 Entity create failed: %s", err)
		return "", common.StringError(err)
	}

	var entity *entityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Unit21Id: %s", entity.Unit21Id)
	return entity.Unit21Id, nil
}

// https://docs.unit21.ai/reference/update_entity
func (e entity) Update(id string, updates any) (err error) {

	_, err = update("entities", id, updates)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return common.StringError(err)
	}

	return nil
}

// https://docs.unit21.ai/reference/add_instruments
func (e entity) AddInstruments(entityId string, instrumentIds []string) (err error) {
	apiKey := os.Getenv("UNIT21_API_KEY")
	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := os.Getenv("UNIT21_URL") + orgName + "/entities/" + entityId + "/add-instruments"

	instruments := make(map[string][]string)
	instruments["instrument_ids"] = instrumentIds
	reqBodyBytes, err := json.Marshal(instruments)
	if err != nil {
		log.Printf("Could not encode instrumentIds to bytes: %s", err)
		return
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for instrumentIds: %s", err)
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("u21-key", apiKey)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed to create instrumentIds: %s", err)
		//handle 409 on update that is not allowed
		//handle 423, 500, 503 for retries
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting return body from instrumentIds add request: %s", err)
		return
	}
	if res.StatusCode != 200 {
		log.Printf("Request failed to create instrumentIds: %s", fmt.Sprint(res.StatusCode))
		err = fmt.Errorf("request failed with status code %s and return body: %s", fmt.Sprint(res.StatusCode), string(body))
		return
	}

	log.Printf("String of body from response: %s", string(body))

	// err = json.Unmarshal(body) // we just need to check if 200
	// if err != nil {
	// 	log.Printf("Reading body failed: %s", err)
	// 	return
	// }

	return
}
