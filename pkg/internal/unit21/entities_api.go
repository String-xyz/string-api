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
	Update(user model.User) (err error)
	AddInstruments(entityId string, instrumentId []string) (err error)
}

type entity struct {
	userRepo    repository.User
	deviceRepo  repository.Device
	contactRepo repository.UserContact
	// platformRepo repository.Platform
}

// With Device and Instrument
// func newEntity(user repository.User, device repository.Device, contact repository.UserContact, instrument repository.Instrument) Entity {
// 	return &entity{userRepo: user, deviceRepo: device, contactRepo: contact, instrumentRepo: instrument}
// }

func newEntity(user repository.User, device repository.Device, contact repository.UserContact) Entity {
	return &entity{userRepo: user, deviceRepo: device, contactRepo: contact}
}

// https://docs.unit21.ai/reference/create_entity
func (e entity) Create(user model.User) (unit21Id string, err error) {

	// ultimately may want a join here.
	// devices, err  := e.deviceRepo.ListByUserId(user.ID, 100, 0)
	// instruments, err  := e.instrumentRepo.ListByUserId(user.ID, 100, 0)

	communications, err := getCommunications(user.ID, e.contactRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity communications: %s", err)
		return "", common.StringError(err)
	}

	digitalData, err := getDigitalData(user.ID, e.deviceRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		return
	}

	// customData, err := getCustomData(user.ID, e.platformRepo)
	// if err != nil {
	// 	log.Printf("Failed to gather Unit21 entity customData: %s", err)
	// 	return "", common.StringError(err)
	// }

	body, err := create("entities", mapUserToEntity(user, communications, digitalData))
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
func (e entity) Update(user model.User) (err error) {

	// ultimately may want a join here.
	// devices, err  := e.deviceRepo.ListByUserId(user.ID, 100, 0)
	// instruments, err  := e.instrumentRepo.ListByUserId(user.ID, 100, 0)

	communications, err := getCommunications(user.ID, e.contactRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity communications: %s", err)
		return "", common.StringError(err)
	}

	digitalData, err := getDigitalData(user.ID, e.deviceRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		return
	}

	// customData, err := getCustomData(user.ID, e.platformRepo)
	// if err != nil {
	// 	log.Printf("Failed to gather Unit21 entity customData: %s", err)
	// 	return "", common.StringError(err)
	// }

	body, err := create("entities", mapUserToEntity(user, communications, digitalData))
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
		return common.StringError(err)
	}

	bodyReader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Printf("Could not create request for instrumentIds: %s", err)
		return common.StringError(err)
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
		return common.StringError(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error extracting return body from instrumentIds add request: %s", err)
		return common.StringError(err)
	}
	if res.StatusCode != 200 {
		log.Printf("Request failed to create instrumentIds: %s", fmt.Sprint(res.StatusCode))
		err = fmt.Errorf("request failed with status code %s and return body: %s", fmt.Sprint(res.StatusCode), string(body))
		return common.StringError(err)
	}

	log.Printf("String of body from response: %s", string(body))

	// err = json.Unmarshal(body) // we just need to check if 200
	// if err != nil {
	// 	log.Printf("Reading body failed: %s", err)
	// 	return
	// }

	return
}

func getCommunications(userId string, contactRepo repository.UserContact) (communications communication, err error) {
	// Get user contacts
	contacts, err := contactRepo.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed go get user contacts: %s", err)
		err = common.StringError(err)
		return
	}

	// Convert contact structs to an entity communication struct
	for _, contact := range contacts {
		if contact.Type == "Email" {
			communications.Emails = append(communications.Emails, contact.Data)
		} else if contact.Type == "Phone" {
			communications.Phones = append(communications.Phones, contact.Data)
		}
	}
	log.Printf("communication: %s", communications)
	return
}

func getDigitalData(userId string, deviceRepo repository.Device) (deviceData digitalData, err error) {
	devices, err := deviceRepo.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed to get user devices: %s", err)
		err = common.StringError(err)
		return
	}

	for _, device := range devices {
		deviceData.IpAddresses = append(deviceData.IpAddresses, device.IpAddresses...)
		deviceData.ClientFingerprints = append(deviceData.ClientFingerprints, device.Fingerprint)
	}
	log.Printf("deviceData: %s", deviceData)
	return
}

func mapUserToEntity(user model.User, communication communication, digitalData digitalData) *u21entity {
	var userTagArr []string
	if user.Tags != nil {
		for key, value := range user.Tags {
			userTagArr = append(userTagArr, key+":"+value)
		}
	}

	jsonBody := &u21entity{
		GeneralData: &general{
			EntityId:     user.ID,
			EntityType:   "user",
			Status:       user.Status,
			RegisteredAt: int(user.CreatedAt.Unix()),
			Tags:         userTagArr, // convert from jsonb into array of key:value string pairs
		},
		UserData: &personal{
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			LastName:   user.LastName,
		},
		CommunicationData: &communication,
		DigitalData:       &digitalData,
		CustomData:        nil, //&custom{
		// 	Platform: nil,//data.partnerName,
		// },
		// add WorkflowOptions if not default?
	}

	return jsonBody
}
