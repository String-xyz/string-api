package unit21

import (
	"context"
	"encoding/json"
	"os"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type Instrument interface {
	Create(ctx context.Context, instrument model.Instrument) (unit21Id string, err error)
	Update(ctx context.Context, instrument model.Instrument) (unit21Id string, err error)
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

func (i instrument) Create(ctx context.Context, instrument model.Instrument) (unit21Id string, err error) {

	source, err := i.getSource(ctx, instrument.UserId)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument source")
		return "", libcommon.StringError(err)
	}

	entities, err := i.getEntities(ctx, instrument.UserId)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument entity")
		return "", libcommon.StringError(err)
	}

	digitalData, err := i.getInstrumentDigitalData(ctx, instrument.UserId)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity digitalData")
		return "", libcommon.StringError(err)
	}

	locationData, err := i.getLocationData(ctx, instrument.LocationId.String)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument location")
		return "", libcommon.StringError(err)
	}

	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/instruments/create"
	body, err := u21Post(url, mapToUnit21Instrument(instrument, source, entities, digitalData, locationData))
	if err != nil {
		log.Err(err).Msg("Unit21 Instrument create failed")
		return "", libcommon.StringError(err)
	}

	var u21Response *createInstrumentResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", libcommon.StringError(err)
	}

	log.Info().Str("Unit21Id", u21Response.Unit21Id).Send()

	// Log create instrument action w/ Unit21
	_, err = i.action.Create(instrument, "Creation", u21Response.Unit21Id, "Creation")
	if err != nil {
		log.Err(err).Msg("Error creating a new instrument action in Unit21")
		return u21Response.Unit21Id, libcommon.StringError(err)
	}

	return u21Response.Unit21Id, nil
}

func (i instrument) Update(ctx context.Context, instrument model.Instrument) (unit21Id string, err error) {

	source, err := i.getSource(ctx, instrument.UserId)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument source")
		return "", libcommon.StringError(err)
	}

	entities, err := i.getEntities(ctx, instrument.UserId)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument entity")
		return "", libcommon.StringError(err)
	}

	digitalData, err := i.getInstrumentDigitalData(ctx, instrument.UserId)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity digitalData")
		return "", libcommon.StringError(err)
	}

	locationData, err := i.getLocationData(ctx, instrument.LocationId.String)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 instrument location")
		return "", libcommon.StringError(err)
	}

	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/" + orgName + "/instruments/" + instrument.Id + "/update"
	body, err := u21Put(url, mapToUnit21Instrument(instrument, source, entities, digitalData, locationData))

	if err != nil {
		log.Err(err).Msg("Unit21 Instrument create failed")
		return "", libcommon.StringError(err)
	}

	var u21Response *updateInstrumentResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", libcommon.StringError(err)
	}

	log.Info().Str("Unit21Id", u21Response.Unit21Id).Send()

	return u21Response.Unit21Id, nil
}

func (i instrument) getSource(ctx context.Context, userId string) (source string, err error) {
	if userId == "" {
		log.Warn().Msg("No userId defined")
		return
	}
	user, err := i.repos.User.GetById(ctx, userId)
	if err != nil {
		log.Err(err).Msg("Failed go get user contacts")
		return "", libcommon.StringError(err)
	}

	if user.Tags["internal"] == "true" {
		return "internal", nil
	}
	return "external", nil
}

func (i instrument) getEntities(ctx context.Context, userId string) (entity instrumentEntity, err error) {
	if userId == "" {
		log.Warn().Msg("No userId defined")
		return
	}

	user, err := i.repos.User.GetById(ctx, userId)
	if err != nil {
		log.Err(err).Msg("Failed go get user contacts")
		err = libcommon.StringError(err)
		return
	}

	entity = instrumentEntity{
		EntityId:       userId,
		RelationshipId: "",
		EntityType:     user.Type,
	}
	return entity, nil
}

func (i instrument) getInstrumentDigitalData(ctx context.Context, userId string) (digitalData instrumentDigitalData, err error) {
	if userId == "" {
		log.Warn().Msg("No userId defined")
		return
	}

	devices, err := i.repos.Device.ListByUserId(ctx, userId, 100, 0)
	if err != nil {
		log.Err(err).Msg("Failed to get user devices")
		err = libcommon.StringError(err)
		return
	}

	for _, device := range devices {
		digitalData.IpAddresses = append(digitalData.IpAddresses, device.IpAddresses...)
	}
	return
}

func (i instrument) getLocationData(ctx context.Context, locationId string) (locationData *instrumentLocationData, err error) {
	if locationId == "" {
		log.Warn().Msg("No locationId defined")
		return
	}

	location, err := i.repos.Location.GetById(ctx, locationId)
	if err != nil {
		log.Err(err).Msg("Failed go get instrument location")
		err = libcommon.StringError(err)
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
		InstrumentId:   instrument.Id,
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
