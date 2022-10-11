package model

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/jmoiron/sqlx/types"
)

// See migrations 0003 for definitions
type TxLeg struct {
	ID           string    `json:"id" db:"id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
	Amount       big.Int   `json:"amount" db:"amount"`
	Value        big.Int   `json:"value" db:"value"`
	AssetID      string    `json:"assetId" db:"asset_id"`
	UserID       string    `json:"userId" db:"user_id"`
	InstrumentID string    `json:"instrumentId" db:"instrument_id"`
}

// See migrations 0003 for definitions
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
	NetworkFee         big.Int        `json:"networkFee" db:"network_fee"`
	ContractParams     types.JSONText `json:"contractParameters" db:"contract_params"`
	ContractFunc       string         `json:"contractFunc" db:"contract_func"`
	TransactionAmount  big.Int        `json:"transactionAmount" db:"transaction_amount"`
	OriginTXLegID      string         `json:"originTXLegId" db:"origin_tx_leg_id"`
	ReceiptTXLegID     string         `json:"receiptTXLegId" db:"receipt_tx_leg_id"`
	ResponseTXLegID    string         `json:"responseTXLegId" db:"response_tx_leg_id"`
	DestinationTXLegID string         `json:"destinationTXLegId" db:"destination_tx_leg_id"`
	ProcessingFee      big.Int        `json:"processingFee" db:"processing_fee"`
	ProcessingFeeAsset string         `json:"processingFeeAsset" db:"processing_fee_asset"`
	StringFee          big.Int        `json:"stringFee" db:"string_fee"`
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
