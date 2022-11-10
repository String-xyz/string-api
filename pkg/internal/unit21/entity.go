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
	Update(user model.User) (unit21Id string, err error)
	AddInstruments(entityId string, instrumentId []string) (err error)
}

type EntityRepos struct {
	Device       repository.Device
	Contact      repository.Contact
	UserPlatform repository.UserPlatform
}

type entity struct {
	repo EntityRepos
}

func NewEntity(r EntityRepos) Entity {
	return &entity{repo: r}
}

// https://docs.unit21.ai/reference/create_entity
func (e entity) Create(user model.User) (unit21Id string, err error) {

	// ultimately may want a join here.

	communications, err := getCommunications(user.ID, e.repo.Contact)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity communications: %s", err)
		return "", common.StringError(err)
	}

	digitalData, err := getEntityDigitalData(user.ID, e.repo.Device)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		return "", common.StringError(err)
	}

	customData, err := getCustomData(user.ID, e.repo.UserPlatform)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity customData: %s", err)
		return "", common.StringError(err)
	}

	body, err := create("entities", mapUserToEntity(user, communications, digitalData, customData))
	if err != nil {
		log.Printf("Unit21 Entity create failed: %s", err)
		return "", common.StringError(err)
	}

	var entity *createEntityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Unit21Id: %s", entity.Unit21Id)

	return entity.Unit21Id, nil
}

// https://docs.unit21.ai/reference/update_entity
func (e entity) Update(user model.User) (unit21Id string, err error) {

	// ultimately may want a join here.

	communications, err := getCommunications(user.ID, e.repo.Contact)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity communications: %s", err)
		err = common.StringError(err)
		return
	}

	digitalData, err := getEntityDigitalData(user.ID, e.repo.Device)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		err = common.StringError(err)
		return
	}

	customData, err := getCustomData(user.ID, e.repo.UserPlatform)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity customData: %s", err)
		err = common.StringError(err)
		return
	}

	body, err := update("entities", user.ID, mapUserToEntity(user, communications, digitalData, customData))
	if err != nil {
		log.Printf("Unit21 Entity create failed: %s", err)
		err = common.StringError(err)
		return
	}

	var entity *updateEntityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		err = common.StringError(err)
		return
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

	return
}

func getCommunications(userId string, contact repository.Contact) (communications entityCommunication, err error) {
	// Get user contacts
	contacts, err := contact.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed to get user contacts: %s", err)
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

func getEntityDigitalData(userId string, device repository.Device) (deviceData entityDigitalData, err error) {
	devices, err := device.ListByUserId(userId, 100, 0)
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

func getCustomData(userId string, userPlatform repository.UserPlatform) (customData entityCustomData, err error) {
	devices, err := userPlatform.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed to get user platforms: %s", err)
		err = common.StringError(err)
		return
	}

	for _, platform := range devices {
		customData.Platforms = append(customData.Platforms, platform.PlatformID)
	}
	log.Printf("deviceData: %s", customData)
	return
}

func mapUserToEntity(user model.User, communication entityCommunication, digitalData entityDigitalData, customData entityCustomData) *u21Entity {
	var userTagArr []string
	if user.Tags != nil {
		for key, value := range user.Tags {
			userTagArr = append(userTagArr, key+":"+value)
		}
	}

	jsonBody := &u21Entity{
		GeneralData: &entityGeneral{
			EntityId:     user.ID,
			EntityType:   "user",
			Status:       user.Status,
			RegisteredAt: int(user.CreatedAt.Unix()),
			Tags:         userTagArr, // convert from jsonb into array of key:value string pairs
		},
		UserData: &entityPersonal{
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			LastName:   user.LastName,
		},
		CommunicationData: &communication,
		DigitalData:       &digitalData,
		CustomData:        &customData,
		// add WorkflowOptions if not default?
	}

	return jsonBody
}
