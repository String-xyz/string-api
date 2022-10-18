package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

// See STRING_USER in Migrations 0001
type User struct {
	ID            string            `json:"id" db:"id"`
	CreatedAt     time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time         `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time        `json:"deactivatedAt" db:"deactivated_at"`
	Type          string            `json:"type" db:"type"`
	Status        string            `json:"status" db:"status"`
	Tags          map[string]string `json:"tags" db:"tags"`
	FirstName     string            `json:"firstName" db:"first_name"`
	MiddleName    string            `json:"middleName" db:"middle_name"`
	LastName      string            `json:"lastName" db:"last_name"`
}

// See PLATFORM in Migrations 0001
type Platform struct {
	ID             string            `json:"id" db:"id"`
	CreatedAt      time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time         `json:"updatedAt" db:"updated_at"`
	DeactivatedAt  *time.Time        `json:"deactivatedAt" db:"deactivated_at"`
	Type           string            `json:"type" db:"type"`
	ApiKey         string            `json:"apiKey" db:"api_key"`
	Authentication AuthType          `json:"authentication" db:"authentication"`
	Tags           map[string]string `json:"Tags" db:"tags"`
}

// See NETWORK in Migrations 0001
type Network struct {
	ID         string    `json:"id" db:"id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	Name       string    `json:"name" db:"name"`
	NetworkID  uint64    `json:"networkId" db:"network_id"`
	ChainID    uint64    `json:"chainId" db:"chain_id"`
	GasTokenID string    `json:"gasTokenId" db:"gas_token_id"`
	GasOracle  string    `json:"gasOracle" db:"gas_oracle"`
	RPCUrl     string    `json:"rpcUrl" db:"rpc_url"`
}

// See ASSET in Migrations 0001
type Asset struct {
	ID          string         `json:"id" db:"id"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	Decimals    uint64         `json:"decimals" db:"decimals"`
	IsCrypto    bool           `json:"isCrypto" db:"is_crypto"`
	NetworkID   sql.NullString `json:"networkId" db:"network_id"`
	ValueOracle sql.NullString `json:"valueOracle" db:"value_oracle"`
}

// See DEVICE in Migrations 0002
type Device struct {
	ID            string         `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
	LastUsedAt    time.Time      `json:"lastUsedAt" db:"last_used_at"`
	ValidatedAt   time.Time      `json:"validatedAt" db:"validated_at"`
	DeactivatedAt time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          string         `json:"type" db:"type"`
	Description   string         `json:"description" db:"description"`
	Fingerprint   string         `json:"fingerprint" db:"fingerprint"`
	IpAddresses   types.JSONText `json:"ipAddresses" db:"ip_addresses"`
	UserID        string         `json:"userId" db:"user_id"`
}

// See CONTACT in Migrations 0002
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

// See LOCATION in Migrations 0002
type Location struct {
	ID             string            `json:"id" db:"id"`
	UserID         string            `json:"userId" db:"user_id"`
	CreatedAt      time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time         `json:"updatedAt" db:"updated_at"`
	Type           string            `json:"type" db:"type"`
	Status         string            `json:"status" db:"status"`
	Tags           map[string]string `json:"tags" db:"tags"`
	BuildingNumber string            `json:"buildingNumber" db:"building_number"`
	UnitNumber     string            `json:"unitNumber" db:"unit_number"`
	StreetName     string            `json:"streetName" db:"street_name"`
	City           string            `json:"city" db:"city"`
	State          string            `json:"state" db:"state"`
	PostalCode     string            `json:"postalCode" db:"postal_code"`
	Country        string            `json:"country" db:"country"`
}

// See INSTRUMENT in Migrations 0002
type Instrument struct {
	ID            string            `json:"id" db:"id"`
	CreatedAt     time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time         `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time        `json:"deactivatedAt" db:"deactivated_at"`
	Type          string            `json:"type" db:"type"`
	Status        string            `json:"status" db:"status"`
	Tags          map[string]string `json:"tags" db:"tags"`
	Network       string            `json:"network" db:"network"`
	PublicKey     string            `json:"publicKey" db:"public_key"`
	Last4         string            `json:"last4" db:"last_4"`
	UserID        string            `json:"userId" db:"user_id"`
	LocationID    string            `json:"locationId" db:"location_id"`
}

// See CONTACT_PLATFORM in Migrations 0003
type ContactPlatform struct {
	ID            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	ContactID     string     `json:"contactId" db:"contact_id"`
	PlatformID    string     `json:"platformId" db:"platform_id"`
}

// See DEVICE_INSTRUMENT in Migrations 0003
type DeviceInstrument struct {
	ID            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	DeviceID      string     `json:"deviceId" db:"device_id"`
	InstrumentID  string     `json:"instrumentId" db:"instrument_id"`
}

// See TX_LEG in Migrations 0003
type TxLeg struct {
	ID           string    `json:"id" db:"id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
	Amount       string    `json:"amount" db:"amount"`
	Value        string    `json:"value" db:"value"`
	AssetID      string    `json:"assetId" db:"asset_id"`
	UserID       string    `json:"userId" db:"user_id"`
	InstrumentID string    `json:"instrumentId" db:"instrument_id"`
}

// See TRANSACTION in Migrations 0003
type Transaction struct {
	ID                 string            `json:"id" db:"id"`
	CreatedAt          time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time         `json:"updatedAt" db:"updated_at"`
	Timestamp          time.Time         `json:"timestamp" db:"timestamp"`
	Type               string            `json:"type" db:"type"`
	Status             string            `json:"status" db:"status"`
	Tags               map[string]string `json:"tags" db:"tags"` // TODO: Fix this alongside Unit21 integration
	DeviceID           string            `json:"deviceId" db:"device_id"`
	IPAddress          string            `json:"ipAddress" db:"ip_address"`
	PlatformID         string            `json:"platformId" db:"platform_id"`
	TransactionHash    string            `json:"transactionHash" db:"transaction_hash"`
	NetworkID          string            `json:"networkId" db:"network_id"`
	NetworkFee         string            `json:"networkFee" db:"network_fee"`
	ContractParams     types.JSONText    `json:"contractParameters" db:"contract_params"`
	ContractFunc       string            `json:"contractFunc" db:"contract_func"`
	TransactionAmount  string            `json:"transactionAmount" db:"transaction_amount"`
	OriginTXLegID      string            `json:"originTXLegId" db:"origin_tx_leg_id"`
	ReceiptTXLegID     string            `json:"receiptTXLegId" db:"receipt_tx_leg_id"`
	ResponseTXLegID    string            `json:"responseTXLegId" db:"response_tx_leg_id"`
	DestinationTXLegID string            `json:"destinationTXLegId" db:"destination_tx_leg_id"`
	ProcessingFee      string            `json:"processingFee" db:"processing_fee"`
	ProcessingFeeAsset string            `json:"processingFeeAsset" db:"processing_fee_asset"`
	StringFee          string            `json:"stringFee" db:"string_fee"`
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
