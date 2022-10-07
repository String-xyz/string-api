package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type UpdateTransaction struct {
	Type               *string         `json:"type" db:"type"`
	Status             *string         `json:"status" db:"status"`
	Tags               *types.JSONText `json:"tags" db:"tags"`
	DeviceID           *string         `json:"deviceId" db:"device_id"`
	IPAddress          *string         `json:"ipAddress" db:"ip_address"`
	PlatformID         *string         `json:"platformId" db:"platform_id"`
	TransactionHash    *string         `json:"transactionHash" db:"transaction_hash"`
	NetworkID          *string         `json:"networkId" db:"network_id"`
	NetworkFee         *uint64         `json:"networkFee" db:"network_fee"` // gas fee
	Parameters         *types.JSONText `json:"parameters" db:"parameters"`
	ContractABI        *string         `json:"contractABI" db:"contract_ABI"`
	OriginTXLegID      *string         `json:"originTXLegId" db:"origin_tx_leg_id"`
	ReceiptTXLegID     *string         `json:"receiptTXLegId" db:"receipt_tx_leg_id"`
	ResponeTXLegID     *string         `json:"responseTXLegId" db:"response_tx_leg_id"`
	DestinationTXLegID *string         `json:"destinationTXLegId" db:"destination_tx_leg_id"`
	ProcessingFee      *float64        `json:"processingFee" db:"processing_fee"`            // GAS IN NATIVE TOKEN
	ProcessingFeeAsset *string         `json:"processingFeeAsset" db:"processing_fee_asset"` // NATIVE TOKEN
	StringFee          *uint64         `json:"stringFee" db:"string_fee"`
}

type UserRegister struct {
	FirstName  string `json:"firstName" db:"first_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type CreatePlatform struct {
	Type           string   `json:"type"`
	Email          string   `json:"email"`
	ApiKey         string   `json:"apiKey" db:"api_key"`
	Authentication AuthType `json:"authentication" db:"authentication"`
}

type UserEmailLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
	FirstNname    *string         `json:"firstName" db:"first_name"`
	MiddleName    *string         `json:"middleName" db:"middle_name"`
	LastName      *string         `json:"lastName" db:"last_name"`
}

type UserContactUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Status        *string    `json:"status" db:"status"`
	Data          *string    `json:"data" db:"data"`
}

type PlaformContactUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Status        *string    `json:"status" db:"status"`
	Data          *string    `json:"data" db:"data"`
}

type UserPKLogin struct {
	PublicAddress string `json:"publicAddress"`
	Signature     string `json:"signature"`
	Nonce         string `json:"nonce"`
}

type EntityType string
type AuthType string
