package model

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Transaction struct {
	ID                 string         `json:"id" db:"id"`
	CreatedAt          time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time      `json:"updatedAt" db:"updated_at"`
	Timestamp          time.Time      `json:"timestamp" db:"timestamp"`
	Type               string         `json:"type" db:"type"`
	Status             string         `json:"status" db:"status"`
	Tags               types.JSONText `json:"tags" db:"tags"`
	DeviceID           string         `json:"deviceId" db:"device_id"`
	IPAddress          string         `json:"ipAddress" db:"ip_address"`
	PlatformID         string         `json:"platformId" db:"platform_id"`
	TransactionHash    string         `json:"transactionHash" db:"transaction_hash"`
	NetworkID          string         `json:"networkId" db:"network_id"`
	NetworkFee         uint64         `json:"networkFee" db:"network_fee"` // gas fee
	Parameters         types.JSONText `json:"parameters" db:"parameters"`
	ContractABI        string         `json:"contractABI" db:"contract_ABI"`
	OriginTXLegID      string         `json:"originTXLegId" db:"origin_tx_leg_id"`
	ReceiptTXLegID     string         `json:"receiptTXLegId" db:"receipt_tx_leg_id"`
	ResponeTXLegID     string         `json:"responseTXLegId" db:"response_tx_leg_id"`
	DestinationTXLegID string         `json:"destinationTXLegId" db:"destination_tx_leg_id"`
	ProcessingFee      float64        `json:"processingFee" db:"processing_fee"`            // GAS IN NATIVE TOKEN
	ProcessingFeeAsset string         `json:"processingFeeAsset" db:"processing_fee_asset"` // NATIVE TOKEN
	StringFee          uint64         `json:"stringFee" db:"string_fee"`
}

type Platform struct {
	ID             string         `json:"id" db:"id"`
	CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt  *time.Time     `json:"deactivatedAt" db:"deactivated_at"`
	Type           string         `json:"type" db:"type"`
	ApiKey         string         `json:"apiKey" db:"api_key"`
	Authentication AuthType       `json:"authentication" db:"authentication"`
	Tags           types.JSONText `json:"Tags" db:"tags"`
}

type User struct {
	ID            string         `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt" db:"deactivated_at"`
	Type          string         `json:"type" db:"type"`
	Status        string         `json:"status" db:"status"`
	Tags          types.JSONText `json:"tags" db:"tags"`
	FirstName     string         `json:"firstName" db:"first_name"`
	MiddleName    string         `json:"middleName" db:"middle_name"`
	LastName      string         `json:"lastName" db:"last_name"`
}

type Contact struct {
	ID                  string     `json:"id" db:"id"`
	UserID              string     `json:"userId" db:"user_id"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
	LastAuthenticatedAt *time.Time `json:"lastAuthenticatedAt" db:"last_authenticated_at"`
	DeactivatedAt       *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type                string     `json:"type" db:"type"`
	Status              string     `json:"status" db:"status"`
	Data                string     `json:"data" db:"data"`
}

type AuthStrategy struct {
	ID            string     `json:"id,omitempty" db:"id"`
	EntityID      string     `json:"entityId" db:"id"` // for redis use only
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"authType" db:"type"`
	EntityType    string     `json:"entityType,omitempty"` // for redis use only
	ContactData   string     `json:"contactData"`          // for redis use only
	ContactID     string     `json:"contactId" db:"contact_id"`
	Data          string     `json:"data" data:"data"`
}

func (a AuthStrategy) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}
