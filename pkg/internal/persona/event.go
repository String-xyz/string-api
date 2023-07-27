package persona

import (
	"encoding/json"
	"time"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/cockroachdb/errors"
)

type EventType string

const (
	EventTypeAccountCreated      = EventType("account.created")
	EventTypeInquiryCreated      = EventType("inquiry.created")
	EventTypeInquiryStarted      = EventType("inquiry.started")
	EventTypeInquiryCompleted    = EventType("inquiry.completed")
	EventTypeVerificationCreated = EventType("verification.created")
	EventTypeVerificationPassed  = EventType("verification.passed")
	EventTypeVerificationFailed  = EventType("verification.failed")
)

type PayloadData interface {
	GetType() string
}

type EventPayload struct {
	Data json.RawMessage `json:"data"`
}

type Event struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes EventAttributes `json:"attributes"`
}

type EventAttributes struct {
	Name      EventType    `json:"name"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Payload   EventPayload `json:"payload"`
}

func (a EventAttributes) GetType() EventType {
	return a.Name
}

func (a EventAttributes) GetPayloadData() (PayloadData, error) {
	switch a.GetType() {
	case EventTypeAccountCreated:
		var data Account
		err := json.Unmarshal(a.Payload.Data, &data)
		if err != nil {
			return nil, common.StringError(err)
		}
		return data, nil
	case EventTypeInquiryCreated, EventTypeInquiryStarted, EventTypeInquiryCompleted:
		var data Inquiry
		err := json.Unmarshal(a.Payload.Data, &data)
		if err != nil {
			return nil, common.StringError(err)
		}
		return data, nil
	case EventTypeVerificationCreated, EventTypeVerificationPassed, EventTypeVerificationFailed:
		var data Verification
		err := json.Unmarshal(a.Payload.Data, &data)
		if err != nil {
			return nil, common.StringError(err)
		}
		return data, nil
	default:
		return nil, common.StringError(errors.Newf("unknown event type: %s", a.GetType()))
	}
}
