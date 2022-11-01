package unit21

import (
	"encoding/json"
	"log"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Instrument interface {
	Create(instrument model.Instrument) (unit21Id string, err error)
	Update(instrument model.Instrument) (unit21Id string, err error)
}

type instrument struct {
	instrumentRepo repository.Instrument
	userRepo       repository.User
	deviceRepo     repository.Device
	locationRepo   repository.Location
}

func NewInstrument(inst repository.Instrument, user repository.User, device repository.Device, location repository.Location) Instrument {
	return &instrument{instrumentRepo: inst, userRepo: user, deviceRepo: device, locationRepo: location}
}

func (i instrument) Create(instrument model.Instrument) (unit21Id string, err error) {

	source, err := getSource(instrument.UserID, i.userRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 instrument source: %s", err)
		return "", common.StringError(err)
	}

	entities, err := getEntities(instrument.UserID, i.userRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 instrument entity: %s", err)
		return "", common.StringError(err)
	}

	digitalData, err := getInstrumentDigitalData(instrument.UserID, i.deviceRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		return "", common.StringError(err)
	}

	locationData, err := getLocationData(instrument.LocationID, i.locationRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 instrument location: %s", err)
		return "", common.StringError(err)
	}

	body, err := create("instruments", mapToUnit21Instrument(instrument, source, entities, digitalData, locationData))
	if err != nil {
		log.Printf("Unit21 Instrument create failed: %s", err)
		return "", common.StringError(err)
	}

	var u21Response *createInstrumentResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Unit21Id: %s", u21Response.Unit21Id)
	return u21Response.Unit21Id, nil
}

func (i instrument) Update(instrument model.Instrument) (unit21Id string, err error) {

	source, err := getSource(instrument.UserID, i.userRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 instrument source: %s", err)
		return "", common.StringError(err)
	}

	entities, err := getEntities(instrument.UserID, i.userRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 instrument entity: %s", err)
		return "", common.StringError(err)
	}

	digitalData, err := getInstrumentDigitalData(instrument.UserID, i.deviceRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 entity digitalData: %s", err)
		return "", common.StringError(err)
	}

	locationData, err := getLocationData(instrument.LocationID, i.locationRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 instrument location: %s", err)
		return "", common.StringError(err)
	}

	body, err := update("instruments", instrument.ID, mapToUnit21Instrument(instrument, source, entities, digitalData, locationData))
	if err != nil {
		log.Printf("Unit21 Instrument create failed: %s", err)
		return "", common.StringError(err)
	}

	var u21Response *updateInstrumentResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Unit21Id: %s", u21Response.Unit21Id)

	return u21Response.Unit21Id, nil
}

func getSource(userID string, userRepo repository.User) (source string, err error) {
	user, err := userRepo.GetById(userID)
	if err != nil {
		log.Printf("Failed go get user contacts: %s", err)
		return "", common.StringError(err)
	}

	if user.Tags["internal"] == "true" {
		return "internal", nil
	}
	return "external", nil
}

func getEntities(userID string, userRepo repository.User) (entity instrumentEntity, err error) {
	user, err := userRepo.GetById(userID)
	if err != nil {
		log.Printf("Failed go get user contacts: %s", err)
		err = common.StringError(err)
		return
	}

	entity = instrumentEntity{
		EntityId:       userID,
		RelationshipId: "",
		EntityType:     user.Type,
	}
	return entity, nil
}

func getInstrumentDigitalData(userId string, deviceRepo repository.Device) (digitalData instrumentDigitalData, err error) {
	devices, err := deviceRepo.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Printf("Failed to get user devices: %s", err)
		err = common.StringError(err)
		return
	}

	for _, device := range devices {
		digitalData.IpAddresses = append(digitalData.IpAddresses, device.IpAddresses...)
	}
	log.Printf("deviceData: %s", digitalData)
	return
}

func getLocationData(locationID string, locationRepo repository.Location) (locationData instrumentLocationData, err error) {
	location, err := locationRepo.GetById(locationID)
	if err != nil {
		log.Printf("Failed go get instrument location: %s", err)
		err = common.StringError(err)
		return
	}

	locationData = instrumentLocationData{
		Type:           location.Type,
		BuildingNumber: location.BuildingNumber,
		UnitNumber:     location.UnitNumber,
		StreetName:     location.StreetName,
		City:           location.City,
		State:          location.State,
		PostalCode:     location.PostalCode,
		Country:        location.Country,
		VerifiedOn:     int(location.CreatedAt.Unix()),
	}

	return locationData, nil
}

func mapToUnit21Instrument(instrument model.Instrument, source string, entityData instrumentEntity, digitalData instrumentDigitalData, locationData instrumentLocationData) *u21Instrument {
	var instrumentTagArr []string
	if instrument.Tags != nil {
		for key, value := range instrument.Tags {
			instrumentTagArr = append(instrumentTagArr, key+":"+value)
		}
	}

	var entityArray []instrumentEntity
	entityArray = append(entityArray, entityData)

	jsonBody := &u21Instrument{
		InstrumentId:   instrument.ID,
		InstrumentType: instrument.Type,
		// InstrumentSubtype:  "",
		// Source:             "internal",
		Status:             instrument.Status,
		RegisteredAt:       int(instrument.CreatedAt.Unix()),
		ParentInstrumentId: "",
		Entities:           entityArray,
		CustomData: &instrumentCustomData{
			None: nil,
		},
		DigitalData:  &digitalData,
		LocationData: &locationData,
		Tags:         instrumentTagArr,
		// Options:      &options,
	}

	log.Printf("%+v\n", jsonBody)

	return jsonBody
}
