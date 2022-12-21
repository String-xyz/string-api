package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

type EntityType string
type AuthType string

type TransactionUpdates struct {
	Type               *string         `json:"type" db:"type"`
	Status             *string         `json:"status" db:"status"`
	Tags               *types.JSONText `json:"tags" db:"tags"`
	DeviceID           *string         `json:"deviceId" db:"device_id"`
	IPAddress          *string         `json:"ipAddress" db:"ip_address"`
	PlatformID         *string         `json:"platformId" db:"platform_id"`
	TransactionHash    *string         `json:"transactionHash" db:"transaction_hash"`
	NetworkID          *string         `json:"networkId" db:"network_id"`
	NetworkFee         *string         `json:"networkFee" db:"network_fee"`
	ContractParams     *pq.StringArray `json:"contractParameters" db:"contract_params"`
	ContractFunc       *string         `json:"contractFunc" db:"contract_func"`
	TransactionAmount  *string         `json:"transactionAmount" db:"transaction_amount"`
	OriginTxLegID      *string         `json:"originTxLegId" db:"origin_tx_leg_id"`
	ReceiptTxLegID     *string         `json:"receiptTxLegId" db:"receipt_tx_leg_id"`
	ResponseTxLegID    *string         `json:"responseTxLegId" db:"response_tx_leg_id"`
	DestinationTxLegID *string         `json:"destinationTxLegId" db:"destination_tx_leg_id"`
	ProcessingFee      *string         `json:"processingFee" db:"processing_fee"`
	ProcessingFeeAsset *string         `json:"processingFeeAsset" db:"processing_fee_asset"`
	StringFee          *string         `json:"stringFee" db:"string_fee"`
}

type UserRegister struct {
	FirstName  string `json:"firstName" db:"first_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
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

type UserPKLogin struct {
	PublicAddress string `json:"publicAddress"`
	Signature     string `json:"signature"`
	Nonce         string `json:"nonce"`
}

// Can be used for user, instrument, etc
type UserRequest struct {
	WalletAddress string `json:"walletAddress" validate:"required"`
	EmailAddress  string `json:"emailAddress" validate:"required,email"`
	FirstName     string `json:"firstName" validate:"required"`
	MiddleName    string `json:"middleName"`
	LastName      string `json:"lastName"`
	Signature     string `json:"signature,omitempty"`
	Password      string `json:"password,omitempty"`
}

type UpdateUserName struct {
	FirstName  string `json:"firstName" db:"first_name" validate:"required"`
	MiddleName string `json:"middleName" db:"middle_name" validate:"required"`
	LastName   string `json:"lastName" db:"last_name" validate:"required"`
}

type ContactUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Status        *string    `json:"status" db:"status"`
	Data          *string    `json:"data" db:"data"`
}

type CreatePlatform struct {
	Type           string   `json:"type"`
	Authentication AuthType `json:"authentication" db:"authentication"`
}

type PlaformContactUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Status        *string    `json:"status" db:"status"`
	Data          *string    `json:"data" db:"data"`
}

type UpdateStatus struct {
	Status *string `json:"status" db:"status"`
}

type NetworkUpdates struct {
	GasTokenID *string `json:"gasTokenId" db:"gas_token_id"`
}

type DeviceUpdates struct {
	ValidatedAt *time.Time `json:"validatedAt" db:"validated_at"`
}
