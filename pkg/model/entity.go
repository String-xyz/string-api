package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// See STRING_USER in Migrations 0001
type User struct {
	Id            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"type" db:"type"`
	Status        string     `json:"status" db:"status"`
	Tags          StringMap  `json:"tags" db:"tags"`
	FirstName     string     `json:"firstName" db:"first_name"`
	MiddleName    string     `json:"middleName" db:"middle_name"`
	LastName      string     `json:"lastName" db:"last_name"`
	Email         string     `json:"email"`
}

// See PLATFORM in Migrations 0005
type Platform struct {
	Id            string         `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ActivatedAt   *time.Time     `json:"activatedAt,omitempty" db:"activated_at"`
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Domains       pq.StringArray `json:"domains" db:"domains"`
	IPAddresses   pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

// See NETWORK in Migrations 0001
type Network struct {
	Id            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string     `json:"name" db:"name"`
	NetworkId     uint64     `json:"networkId" db:"network_id"`
	ChainId       uint64     `json:"chainId" db:"chain_id"`
	GasTokenId    string     `json:"gasTokenId" db:"gas_token_id"`
	GasOracle     string     `json:"gasOracle" db:"gas_oracle"`
	RPCUrl        string     `json:"rpcUrl" db:"rpc_url"`
	ExplorerUrl   string     `json:"explorerUrl" db:"explorer_url"`
}

// See ASSET in Migrations 0001
type Asset struct {
	Id            string         `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Decimals      uint64         `json:"decimals" db:"decimals"`
	IsCrypto      bool           `json:"isCrypto" db:"is_crypto"`
	NetworkId     sql.NullString `json:"networkId" db:"network_id"`
	ValueOracle   sql.NullString `json:"valueOracle" db:"value_oracle"`
	ValueOracle2  sql.NullString `json:"valueOracle2" db:"value_oracle_2"`
}

// See USER_PLATFORM in Migrations 0002
type UserToPlatform struct {
	UserId     string `json:"userId" db:"user_id"`
	PlatformId string `json:"platformId" db:"platform_id"`
}

// See DEVICE in Migrations 0002
type Device struct {
	Id            string         `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
	LastUsedAt    time.Time      `json:"lastUsedAt" db:"last_used_at"`
	ValidatedAt   *time.Time     `json:"validatedAt" db:"validated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string         `json:"type" db:"type"`
	Description   string         `json:"description" db:"description"`
	Fingerprint   string         `json:"fingerprint" db:"fingerprint"`
	IpAddresses   pq.StringArray `json:"ipAddresses,omitempty" db:"ip_addresses"`
	UserId        string         `json:"userId" db:"user_id"`
}

// See CONTACT in Migrations 0002
type Contact struct {
	Id                  string     `json:"id" db:"id"`
	UserId              string     `json:"userId" db:"user_id"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
	LastAuthenticatedAt *time.Time `json:"lastAuthenticatedAt" db:"last_authenticated_at"`
	ValidatedAt         *time.Time `json:"validatedAt" db:"validated_at"`
	DeactivatedAt       *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type                string     `json:"type" db:"type"`
	Status              string     `json:"status" db:"status"`
	Data                string     `json:"data" db:"data"`
}

// See LOCATION in Migrations 0002
type Location struct {
	Id             string     `json:"id" db:"id"`
	UserId         string     `json:"userId" db:"user_id"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt  *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type           string     `json:"type" db:"type"`
	Status         string     `json:"status" db:"status"`
	Tags           StringMap  `json:"tags" db:"tags"`
	BuildingNumber string     `json:"buildingNumber" db:"building_number"`
	UnitNumber     string     `json:"unitNumber" db:"unit_number"`
	StreetName     string     `json:"streetName" db:"street_name"`
	City           string     `json:"city" db:"city"`
	State          string     `json:"state" db:"state"`
	PostalCode     string     `json:"postalCode" db:"postal_code"`
	Country        string     `json:"country" db:"country"`
}

// See INSTRUMENT in Migrations 0002
type Instrument struct {
	Id            string         `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string         `json:"type" db:"type"`
	Status        string         `json:"status" db:"status"`
	Tags          StringMap      `json:"tags" db:"tags"`
	Network       string         `json:"network" db:"network"`
	PublicKey     string         `json:"publicKey" db:"public_key"`
	Last4         string         `json:"last4" db:"last_4"`
	UserId        string         `json:"userId" db:"user_id"`
	LocationId    sql.NullString `json:"locationId" db:"location_id"`
	Name          string         `json:"name" db:"name"`
}

// See CONTACT_PLATFORM in Migrations 0003
type ContactToPlatform struct {
	ContactId  string `json:"contactId" db:"contact_id"`
	PlatformId string `json:"platformId" db:"platform_id"`
}

// See DEVICE_INSTRUMENT in Migrations 0003
type DeviceToInstrument struct {
	DeviceId     string `json:"deviceId" db:"device_id"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
}

// See Tx_LEG in Migrations 0003
type TxLeg struct {
	Id            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Timestamp     time.Time  `json:"timestamp" db:"timestamp"`
	Amount        string     `json:"amount" db:"amount"`
	Value         string     `json:"value" db:"value"`
	AssetId       string     `json:"assetId" db:"asset_id"`
	UserId        string     `json:"userId" db:"user_id"`
	InstrumentId  string     `json:"instrumentId" db:"instrument_id"`
}

// See TRANSACTION in Migrations 0003
type Transaction struct {
	Id                 string         `json:"id" db:"id"`
	CreatedAt          time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt      *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type               string         `json:"type,omitempty" db:"type"`
	Status             string         `json:"status,omitempty" db:"status"`
	Tags               StringMap      `json:"tags,omitempty" db:"tags"`
	DeviceId           string         `json:"deviceId,omitempty" db:"device_id"`
	IPAddress          string         `json:"ipAddress,omitempty" db:"ip_address"`
	PlatformId         string         `json:"platformId,omitempty" db:"platform_id"`
	TransactionHash    string         `json:"transactionHash,omitempty" db:"transaction_hash"`
	NetworkId          string         `json:"networkId,omitempty" db:"network_id"`
	NetworkFee         string         `json:"networkFee,omitempty" db:"network_fee"`
	ContractParams     pq.StringArray `json:"contractParameters,omitempty" db:"contract_params"`
	ContractFunc       string         `json:"contractFunc,omitempty" db:"contract_func"`
	TransactionAmount  string         `json:"transactionAmount,omitempty" db:"transaction_amount"`
	OriginTxLegId      string         `json:"originTxLegId,omitempty" db:"origin_tx_leg_id"`
	ReceiptTxLegId     sql.NullString `json:"receiptTxLegId,omitempty" db:"receipt_tx_leg_id"`
	ResponseTxLegId    sql.NullString `json:"responseTxLegId,omitempty" db:"response_tx_leg_id"`
	DestinationTxLegId string         `json:"destinationTxLegId,omitempty" db:"destination_tx_leg_id"`
	ProcessingFee      string         `json:"processingFee,omitempty" db:"processing_fee"`
	ProcessingFeeAsset string         `json:"processingFeeAsset,omitempty" db:"processing_fee_asset"`
	StringFee          string         `json:"stringFee,omitempty" db:"string_fee"`
	PaymentCode        string         `json:"paymentCode,omitempty" db:"payment_code"`
}

type AuthStrategy struct {
	Id            string         `json:"id,omitempty" db:"id"`
	Status        string         `json:"status" db:"status"`
	EntityId      string         `json:"entityId,omitempty"` // for redis use only
	Type          string         `json:"authType" db:"type"`
	EntityType    string         `json:"entityType,omitempty"`  // for redis use only
	ContactData   string         `json:"contactData,omitempty"` // for redis use only
	ContactId     NullableString `json:"contactId,omitempty" db:"contact_id"`
	Data          string         `json:"data" data:"data"`
	CreatedAt     time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	ExpiresAt     time.Time      `json:"expireAt,omitempty" db:"expire_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
}

type Apikey struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"type" db:"type"`
	Public        string     `json:"public" db:"public"` // an unhased public key
	Secret        *string    `json:"secret" db:"secret"` // a hashed secret key
	Description   *string    `json:"description" db:"description"`
	CreatedBy     string     `json:"createdBy" db:"created_by"`
	PlatformId    string     `json:"platformId" db:"platform_id"`
}

type Contract struct {
	ID            string         `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string         `json:"name" db:"name"`
	Address       string         `json:"address" db:"address"`
	Functions     pq.StringArray `json:"functions" db:"functions"`
	NetworkID     string         `json:"networkId" db:"network_id"`
	PlatformID    string         `json:"platformId" db:"platform_id"`
}

func (a AuthStrategy) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}
