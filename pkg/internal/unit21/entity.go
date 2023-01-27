package unit21

import (
	"encoding/json"
	"log"
	"os"

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
	Device         repository.Device
	Contact        repository.Contact
	UserToPlatform repository.UserToPlatform
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

	communications, err := e.getCommunications(user.ID)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity communications: %s", err)
		return "", common.StringError(err)
	}

	digitalData, err := e.getEntityDigitalData(user.ID)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		return "", common.StringError(err)
	}

	customData, err := e.getCustomData(user.ID)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity customData: %s", err)
		return "", common.StringError(err)
	}

	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/entities/create"
	body, err := u21Post(url, mapUserToEntity(user, communications, digitalData, customData))
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

	communications, err := e.getCommunications(user.ID)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity communications: %s", err)
		err = common.StringError(err)
		return
	}

	digitalData, err := e.getEntityDigitalData(user.ID)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		err = common.StringError(err)
		return
	}

	customData, err := e.getCustomData(user.ID)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity customData: %s", err)
		err = common.StringError(err)
		return
	}

	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/" + orgName + "/entities/" + user.ID + "/update"
	body, err := u21Put(url, mapUserToEntity(user, communications, digitalData, customData))

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
	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/" + orgName + "/entities/" + entityId + "/add-instruments"

	instruments := make(map[string][]string)
	instruments["instrument_ids"] = instrumentIds

	body, err := u21Put(url, instruments)
	if err != nil {
		log.Printf("Unit21 Entity Add Instruments failed: %s", err)
		err = common.StringError(err)
		return
	}

	log.Printf("String of body from response: %s", string(body))

	return
}

func (e entity) getCommunications(userId string) (communications entityCommunication, err error) {
	// Get user contacts
	contacts, err := e.repo.Contact.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed to get user contacts: %s", err)
		err = common.StringError(err)
		return
	}

	// Convert contact structs to an entity communication struct
	for _, contact := range contacts {
		if contact.Type == "email" {
			communications.Emails = append(communications.Emails, contact.Data)
		} else if contact.Type == "phone" {
			communications.Phones = append(communications.Phones, contact.Data)
		}
	}
	return
}

func (e entity) getEntityDigitalData(userId string) (deviceData entityDigitalData, err error) {
	devices, err := e.repo.Device.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed to get user devices: %s", err)
		err = common.StringError(err)
		return
	}

	for _, device := range devices {
		deviceData.IpAddresses = append(deviceData.IpAddresses, device.IpAddresses...)
		deviceData.ClientFingerprints = append(deviceData.ClientFingerprints, device.Fingerprint)
	}
	return
}

func (e entity) getCustomData(userId string) (customData entityCustomData, err error) {
	devices, err := e.repo.UserToPlatform.ListByUserId(userId, 100, 0)
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
