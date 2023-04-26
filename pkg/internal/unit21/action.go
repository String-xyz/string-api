package unit21

import (
	"encoding/json"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/env"
	"github.com/String-xyz/string-api/pkg/internal/common"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/rs/zerolog/log"
)

type Action interface {
	Create(instrument model.Instrument,
		actionDetails string,
		unit21InstrumentId string,
		eventSubtype string) (unit21Id string, err error)
}

type action struct {
}

func NewAction() Action {
	return &action{}
}

func (a action) Create(
	instrument model.Instrument,
	actionDetails string,
	unit21InstrumentId string,
	eventSubtype string) (unit21Id string, err error) {

	actionData := actionData{
		ActionType:    instrument.Type,
		ActionDetails: actionDetails,
		EntityId:      instrument.UserId,
		EntityType:    "user",
		InstrumentId:  instrument.Id,
	}

	url := "https://" + env.Var.UNIT21_ENV + ".unit21.com/v1/events/create"
	body, err := u21Post(url, mapToUnit21ActionEvent(instrument, actionData, unit21InstrumentId, eventSubtype))
	if err != nil {
		log.Err(err).Msg("Unit21 Action create failed")
		return "", libcommon.StringError(err)
	}

	var u21Response *createEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", libcommon.StringError(err)
	}

	log.Info().Str("unit21Id", u21Response.Unit21Id).Msg("Create Action")

	return u21Response.Unit21Id, nil
}

func mapToUnit21ActionEvent(instrument model.Instrument, actionData actionData, unit21Id string, eventSubtype string) *u21Event {
	var instrumentTagArr []string
	if instrument.Tags != nil {
		for key, value := range instrument.Tags {
			instrumentTagArr = append(instrumentTagArr, key+":"+value)
		}
	}

	jsonBody := &u21Event{
		GeneralData: &eventGeneral{
			EventId:      unit21Id,                         //required
			EventType:    "action",                         //required
			EventTime:    int(instrument.CreatedAt.Unix()), //required
			EventSubtype: eventSubtype,                     //required for RTR
			Status:       instrument.Status,
			Parents:      nil,
			Tags:         instrumentTagArr,
		},
		ActionData:   &actionData,
		DigitalData:  nil,
		LocationData: nil,
		CustomData:   nil,
	}

	actionBody, err := common.BetterStringify(jsonBody)
	if err != nil {
		log.Err(err).Msg("Error creating action body")
		return jsonBody
	}
	log.Info().Str("body", actionBody).Msg("Create Action action body")

	return jsonBody
}
