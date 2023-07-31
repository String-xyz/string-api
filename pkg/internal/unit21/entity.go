package unit21

import (
	"context"
	"encoding/json"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type Entity interface {
	Create(ctx context.Context, user model.User) (unit21Id string, err error)
	Update(ctx context.Context, user model.User) (unit21Id string, err error)
	AddInstruments(entityId string, instrumentId []string) (err error)
}

type EntityRepos struct {
	Device  repository.Device
	Contact repository.Contact
	User    repository.User
}

type entity struct {
	repo EntityRepos
}

func NewEntity(r EntityRepos) Entity {
	return &entity{repo: r}
}

// https://docs.unit21.ai/reference/create_entity
func (e entity) Create(ctx context.Context, user model.User) (unit21Id string, err error) {
	// ultimately may want a join here.

	communications, err := e.getCommunications(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity communications")
		return "", libcommon.StringError(err)
	}

	digitalData, err := e.getEntityDigitalData(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity digitalData")
		return "", libcommon.StringError(err)
	}

	customData, err := e.getCustomData(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity customData")
		return "", libcommon.StringError(err)
	}

	url := "https://" + config.Var.UNIT21_ENV + ".unit21.com/v1/entities/create"
	body, err := u21Post(url, mapUserToEntity(user, communications, digitalData, customData))
	if err != nil {
		log.Err(err).Msg("Unit21 Entity create failed")
		return "", libcommon.StringError(err)
	}

	var entity *createEntityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", libcommon.StringError(err)
	}

	log.Info().Str("Unit21Id", entity.Unit21Id).Send()

	return entity.Unit21Id, nil
}

// https://docs.unit21.ai/reference/update_entity
func (e entity) Update(ctx context.Context, user model.User) (unit21Id string, err error) {

	// ultimately may want a join here.

	communications, err := e.getCommunications(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity communications")
		err = libcommon.StringError(err)
		return
	}

	digitalData, err := e.getEntityDigitalData(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity digitalData")
		err = libcommon.StringError(err)
		return
	}

	customData, err := e.getCustomData(ctx, user.Id)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 entity customData")
		err = libcommon.StringError(err)
		return
	}

	orgName := config.Var.UNIT21_ORG_NAME
	url := "https://" + config.Var.UNIT21_ENV + ".unit21.com/v1/" + orgName + "/entities/" + user.Id + "/update"
	body, err := u21Put(url, mapUserToEntity(user, communications, digitalData, customData))

	if err != nil {
		log.Err(err).Msg("Unit21 Entity create failed")
		err = libcommon.StringError(err)
		return
	}

	var entity *updateEntityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		err = libcommon.StringError(err)
		return
	}

	log.Info().Str("Unit21Id", entity.Unit21Id).Send()
	return entity.Unit21Id, nil
}

// https://docs.unit21.ai/reference/add_instruments
func (e entity) AddInstruments(entityId string, instrumentIds []string) (err error) {
	orgName := config.Var.UNIT21_ORG_NAME
	url := "https://" + config.Var.UNIT21_ENV + ".unit21.com/v1/" + orgName + "/entities/" + entityId + "/add-instruments"

	instruments := make(map[string][]string)
	instruments["instrument_ids"] = instrumentIds

	_, err = u21Put(url, instruments)
	if err != nil {
		log.Err(err).Msg("Unit21 Entity Add Instruments failed")
		err = libcommon.StringError(err)
		return
	}

	return
}

func (e entity) getCommunications(ctx context.Context, userId string) (communications entityCommunication, err error) {
	// Get user contacts
	contacts, err := e.repo.Contact.ListByUserId(ctx, userId, 100, 0)
	if err != nil {
		log.Err(err).Msg("Failed to get user contacts")
		err = libcommon.StringError(err)
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

func (e entity) getEntityDigitalData(ctx context.Context, userId string) (deviceData entityDigitalData, err error) {
	devices, err := e.repo.Device.ListByUserId(ctx, userId, 100, 0)
	if err != nil {
		log.Err(err).Msg("Failed to get user devices")
		err = libcommon.StringError(err)
		return
	}

	for _, device := range devices {
		deviceData.IpAddresses = append(deviceData.IpAddresses, device.IpAddresses...)
		deviceData.ClientFingerprints = append(deviceData.ClientFingerprints, device.Fingerprint)
	}
	return
}

func (e entity) getCustomData(ctx context.Context, userId string) (customData entityCustomData, err error) {
	platforms, err := e.repo.User.GetPlatforms(ctx, userId, 100, 0)
	if err != nil {
		log.Err(err).Msg("Failed to get user platforms")
		err = libcommon.StringError(err)
		return
	}

	for _, platform := range platforms {
		customData.Platforms = append(customData.Platforms, platform.Id)
	}
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
			EntityId:     user.Id,
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
