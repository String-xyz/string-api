package unit21

import (
	"encoding/json"
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type Instrument interface {
	Create(instrument model.Instrument) (unit21Id string, err error)
	Update(instrument model.Instrument) (unit21Id string, err error)
}

type InstrumentRepos struct {
	User     repository.User
	Device   repository.Device
	Location repository.Location
}

type instrument struct {
	action Action
	repos  InstrumentRepos
}

func NewInstrument(r InstrumentRepos, a Action) Instrument {
	return &instrument{repos: r, action: a}
}

func (i instrument) Create(instrument model.Instrument) (unit21Id string, err error) {

	source, err := i.getSource(instrument.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument source")
		return "", common.StringError(err)
	}

	entities, err := i.getEntities(instrument.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument entity")
		return "", common.StringError(err)
	}

	digitalData, err := i.getInstrumentDigitalData(instrument.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity digitalData")
		return "", common.StringError(err)
	}

	locationData, err := i.getLocationData(instrument.LocationID.String)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument location")
		return "", common.StringError(err)
	}

	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/instruments/create"
	body, err := u21Post(url, mapToUnit21Instrument(instrument, source, entities, digitalData, locationData))
	if err != nil {
		log.Err(err).Msg("Unit21 Instrument create failed")
		return "", common.StringError(err)
	}

	var u21Response *createInstrumentResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", common.StringError(err)
	}

	log.Info().Str("Unit21Id", u21Response.Unit21Id).Send()

	// Log create instrument action w/ Unit21
	_, err = i.action.Create(instrument, "Creation", u21Response.Unit21Id, "Creation")
	if err != nil {
		log.Err(err).Msg("Error creating a new instrument action in Unit21")
		return u21Response.Unit21Id, common.StringError(err)
	}

	return u21Response.Unit21Id, nil
}

func (i instrument) Update(instrument model.Instrument) (unit21Id string, err error) {

	source, err := i.getSource(instrument.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument source")
		return "", common.StringError(err)
	}

	entities, err := i.getEntities(instrument.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument entity")
		return "", common.StringError(err)
	}

	digitalData, err := i.getInstrumentDigitalData(instrument.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity digitalData")
		return "", common.StringError(err)
	}

	locationData, err := i.getLocationData(instrument.LocationID.String)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument location")
		return "", common.StringError(err)
	}

	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/" + orgName + "/instruments/" + instrument.ID + "/update"
	body, err := u21Put(url, mapToUnit21Instrument(instrument, source, entities, digitalData, locationData))

	if err != nil {
		log.Err(err).Msg("Unit21 Instrument create failed")
		return "", common.StringError(err)
	}

	var u21Response *updateInstrumentResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", common.StringError(err)
	}

	log.Info().Str("Unit21Id", u21Response.Unit21Id).Send()

	return u21Response.Unit21Id, nil
}

func (i instrument) getSource(userId string) (source string, err error) {
	if userId == "" {
		log.Warn().Msg("No userId defined")
		return
	}
	user, err := i.repos.User.GetById(userId)
	if err != nil {
		log.Err(err).Msg("Failed go get user contacts")
		return "", common.StringError(err)
	}

	if user.Tags["internal"] == "true" {
		return "internal", nil
	}
	return "external", nil
}

func (i instrument) getEntities(userId string) (entity instrumentEntity, err error) {
	if userId == "" {
		log.Warn().Msg("No userId defined")
		return
	}

	user, err := i.repos.User.GetById(userId)
	if err != nil {
		log.Err(err).Msg("Failed go get user contacts")
		err = common.StringError(err)
		return
	}

	entity = instrumentEntity{
		EntityId:       userId,
		RelationshipId: "",
		EntityType:     user.Type,
	}
	return entity, nil
}

func (i instrument) getInstrumentDigitalData(userId string) (digitalData instrumentDigitalData, err error) {
	if userId == "" {
		log.Warn().Msg("No userId defined")
		return
	}

	devices, err := i.repos.Device.ListByUserId(userId, 100, 0)
	if err != nil {
		log.Err(err).Msg("Failed to get user devices")
		err = common.StringError(err)
		return
	}

	for _, device := range devices {
		digitalData.IpAddresses = append(digitalData.IpAddresses, device.IpAddresses...)
	}
	return
}

func (i instrument) getLocationData(locationId string) (locationData *instrumentLocationData, err error) {
	if locationId == "" {
		log.Warn().Msg("No locationId defined")
		return
	}

	location, err := i.repos.Location.GetById(locationId)
	if err != nil {
		log.Err(err).Msg("Failed go get instrument location")
		err = common.StringError(err)
		return
	}
	if location.CreatedAt.Unix() != 0 {
		locationData = &instrumentLocationData{
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
	}

	return locationData, nil
}

func mapToUnit21Instrument(instrument model.Instrument, source string, entityData instrumentEntity, digitalData instrumentDigitalData, locationData *instrumentLocationData) *u21Instrument {
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
		CustomData:         nil, //TODO: include platform in customData
		DigitalData:        &digitalData,
		LocationData:       locationData,
		Tags:               instrumentTagArr,
		// Options:      &options,
	}

	return jsonBody
}
