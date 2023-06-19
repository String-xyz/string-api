package checkout

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	AuthorizationApprovedEvent EventType = "authorization_approved"
	AthorizationDeclinedEvent  EventType = "authorization_declined"
	PaymentApprovedEvent       EventType = "payment_approved"
	PaymentCapturedEvent       EventType = "Payment_captured"
	PaymentDeclinedEvent       EventType = "payment_declined"
	PaymentRefundedEvent       EventType = "payment_refunded"
	PaymentaReturedEvent       EventType = "payment_returned"
	PaymentPendingEvent        EventType = "payment_pending"
	PaymentVoidedEvent         EventType = "payment_voided"
)

type WebhookEvent struct {
	Id        string       `json:"id,omitempty"`
	Type      EventType    `json:"type,omitempty"`
	CreatedOn string       `json:"created_on,omitempty"`
	Data      EventPayload `json:"data,omitempty"`
}

type AuthorizationEvent struct {
	CardId               string    `json:"card_id,omitempty"`
	TransactionId        string    `json:"transaction_id,omitempty"`
	TransactionType      string    `json:"transaction_type,omitempty"`
	TransmissionDateTime time.Time `json:"transmission_date_time,omitempty"`
	AuthorizationType    string    `json:"authorization_type,omitempty"`
	TransactionAmount    int       `json:"transaction_amount,omitempty"`
	TransactionCurrency  string    `json:"transaction_currency,omitempty"`
	BillingAmount        int       `json:"billing_amount,omitempty"`
}

type PaymentEvent struct {
	Id              string    `json:"id,omitempty"`
	ActionId        string    `json:"action_id,omitempty"`
	Reference       string    `json:"reference,omitempty"`
	Amount          int       `json:"amount,omitempty"`
	AuthCode        string    `json:"auth_code,omitempty"`
	Currency        string    `json:"currency,omitempty"`
	PaymentType     string    `json:"payment_type,omitempty"`
	ProccesedOn     time.Time `json:"processed_on,omitempty"`
	ResponseCode    string    `json:"response_code,omitempty"`
	ResponseSummary string    `json:"response_summary,omitempty"`
}

type AuthorizationApproved AuthorizationEvent

type AuthorizationDeclined struct {
	AuthorizationEvent
	DeclineReason string `json:"decline_reason,omitempty"`
}

type (
	PaymentApproved PaymentEvent
	PaymentCaptured PaymentEvent
	PaymentDeclined PaymentEvent
	PaymentRefunded PaymentEvent
	PaymentReturned PaymentEvent
	PaymentPending  PaymentEvent
	PaymentVoided   PaymentEvent
)

type EventPayload interface {
	GetType() EventType
}

func (a AuthorizationApproved) GetType() EventType {
	return AuthorizationApprovedEvent
}

func (a AuthorizationDeclined) GetType() EventType {
	return AthorizationDeclinedEvent
}

func (a PaymentApproved) GetType() EventType {
	return PaymentApprovedEvent
}

func (a PaymentCaptured) GetType() EventType {
	return PaymentCapturedEvent
}

func (a PaymentDeclined) GetType() EventType {
	return PaymentDeclinedEvent
}

// for the time being only 4 events are supported
// if we need to support more events we need to add them here
func (e *WebhookEvent) UnmarshalJSON(data []byte) error {
	type Alias WebhookEvent
	alias := struct {
		*Alias
		// RawData lets us delay parsing the data field until we know the type
		RawData json.RawMessage `json:"data,omitempty"`
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	switch e.Type {
	case AuthorizationApprovedEvent:
		var a AuthorizationApproved
		err := json.Unmarshal(alias.RawData, &a)
		e.Data = a
		return err
	case AthorizationDeclinedEvent:
		var a AuthorizationDeclined
		err := json.Unmarshal(alias.RawData, &a)
		e.Data = a
		return err
	case PaymentApprovedEvent:
		var p PaymentApproved
		err := json.Unmarshal(alias.RawData, &p)
		e.Data = p
		return err
	case PaymentCapturedEvent:
		var p PaymentCaptured
		err := json.Unmarshal(alias.RawData, &p)
		e.Data = p
		return err
	}

	return nil
}
