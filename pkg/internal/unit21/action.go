package unit21

import (
	"encoding/json"
	"log"
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Action interface {
	Create(instrument model.Instrument,
		actionType string,
		actionDetails string,
		unit21InstrumentId string,
		eventSubtype string) (unit21Id string, err error)
}

type ActionRepo struct {
	User     repository.User
	Device   repository.Device
	Location repository.Location
}

type action struct {
	repo ActionRepo
}

func NewAction(r ActionRepo) Action {
	return &action{repo: r}
}

func (a action) Create(
	instrument model.Instrument,
	actionType string,
	actionDetails string,
	unit21InstrumentId string,
	eventSubtype string) (unit21Id string, err error) {

	actionData := actionData{
		ActionType:    actionType,
		ActionDetails: actionDetails,
		EntityId:      instrument.UserID,
		EntityType:    "user",
		InstrumentId:  instrument.ID,
	}

	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/events/create"
	body, err := u21Post(url, mapToUnit21ActionEvent(instrument, actionData, unit21InstrumentId, eventSubtype))
	if err != nil {
		log.Printf("Unit21 Action create failed: %s", err)
		return "", common.StringError(err)
	}

	var u21Response *createEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Create Action Unit21Id: %s", u21Response.Unit21Id)

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
		log.Printf("\nError creating action body\n")
		return jsonBody
	}
	log.Printf("\nCreate Action action body: %+v\n", actionBody)

	return jsonBody
}
