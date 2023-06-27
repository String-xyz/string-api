package model

import (
	"database/sql"
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
	DeviceId           *string         `json:"deviceId" db:"device_id"`
	IPAddress          *string         `json:"ipAddress" db:"ip_address"`
	PlatformId         *string         `json:"platformId" db:"platform_id"`
	TransactionHash    *string         `json:"transactionHash" db:"transaction_hash"`
	NetworkId          *string         `json:"networkId" db:"network_id"`
	NetworkFee         *string         `json:"networkFee" db:"network_fee"`
	ContractParams     *pq.StringArray `json:"contractParameters" db:"contract_params"`
	ContractFunc       *string         `json:"contractFunc" db:"contract_func"`
	TransactionAmount  *string         `json:"transactionAmount" db:"transaction_amount"`
	OriginTxLegId      *string         `json:"originTxLegId" db:"origin_tx_leg_id"`
	ReceiptTxLegId     *string         `json:"receiptTxLegId" db:"receipt_tx_leg_id"`
	ResponseTxLegId    *string         `json:"responseTxLegId" db:"response_tx_leg_id"`
	DestinationTxLegId *string         `json:"destinationTxLegId" db:"destination_tx_leg_id"`
	ProcessingFee      *string         `json:"processingFee" db:"processing_fee"`
	ProcessingFeeAsset *string         `json:"processingFeeAsset" db:"processing_fee_asset"`
	StringFee          *string         `json:"stringFee" db:"string_fee"`
	PaymentCode        *string         `json:"paymentCode" db:"payment_code"`
}

type InstrumentUpdates struct {
	Type       *string         `json:"type" db:"type"`
	Status     *string         `json:"status" db:"status"`
	Tags       *StringMap      `json:"tags" db:"tags"`
	Network    *string         `json:"network" db:"network"`
	PublicKey  *string         `json:"publicKey" db:"public_key"`
	Last4      *string         `json:"last4" db:"last_4"`
	UserId     *string         `json:"userId" db:"user_id"`
	LocationId *sql.NullString `json:"locationId" db:"location_id"`
}

type TxLegUpdates struct {
	Timestamp    *time.Time `json:"timestamp" db:"timestamp"`
	Amount       *string    `json:"amount" db:"amount"`
	Value        *string    `json:"value" db:"value"`
	AssetId      *string    `json:"assetId" db:"asset_id"`
	UserId       *string    `json:"userId" db:"user_id"`
	InstrumentId *string    `json:"instrumentId" db:"instrument_id"`
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
	Type       *string         `json:"type" db:"type"`
	Status     *string         `json:"status" db:"status"`
	CheckoutId *string         `json:"checkoutId" db:"checkout_id"`
	Tags       *types.JSONText `json:"tags" db:"tags"`
	FirstName  *string         `json:"firstName" db:"first_name"`
	MiddleName *string         `json:"middleName" db:"middle_name"`
	LastName   *string         `json:"lastName" db:"last_name"`
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
	FirstName  string `json:"firstName" db:"first_name" validate:"max=255"`
	MiddleName string `json:"middleName" db:"middle_name" validate:"max=255"`
	LastName   string `json:"lastName" db:"last_name" validate:"max=255"`
}

type ContactUpdates struct {
	Type   *string `json:"type" db:"type"`
	Status *string `json:"status" db:"status"`
	Data   *string `json:"data" db:"data"`
}

type CreatePlatform struct {
	Type           string   `json:"type"`
	Authentication AuthType `json:"authentication" db:"authentication"`
}

type PlatformContactUpdates struct {
	Type   *string `json:"type" db:"type"`
	Status *string `json:"status" db:"status"`
	Data   *string `json:"data" db:"data"`
}

type UpdateStatus struct {
	Status *string `json:"status" db:"status"`
}

type NetworkUpdates struct {
	GasTokenId *string `json:"gasTokenId" db:"gas_token_id"`
}

type DeviceUpdates struct {
	ValidatedAt *time.Time      `json:"validatedAt" db:"validated_at"`
	IpAddresses *pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type RefreshTokenPayload struct {
	WalletAddress string `json:"walletAddress" validate:"required,eth_addr"`
}

type PreValidateEmail struct {
	Email string `json:"email" validate:"required,email"`
}
