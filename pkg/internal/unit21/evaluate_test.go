package unit21

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// This transaction should pass
func TestEvaluateTransactionPass(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	transactionId := uuid.NewString()
	OriginTxLegID := uuid.NewString()
	DestinationTxLegID := uuid.NewString()
	userId1 := uuid.NewString()
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	networkId := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          networkId,
		NetworkFee:         "100000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "1000000000",
		OriginTxLegID:      OriginTxLegID,
		ReceiptTxLegID:     sql.NullString{String: uuid.NewString()},
		ResponseTxLegID:    sql.NullString{String: uuid.NewString()},
		DestinationTxLegID: DestinationTxLegID,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(OriginTxLegID, time.Now(), "1000000", "1000000", assetId1, userId1, instrumentId1)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(OriginTxLegID).WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(DestinationTxLegID, time.Now(), "1", "10000000", assetId2, userId1, instrumentId2)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(DestinationTxLegID).WillReturnRows(mockedTxLegRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId1, "USD", "fiat USD", 6, false, networkId, "self")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId1).WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId2, "Noose The Goose", "Noose the Goose NFT", 0, true, networkId, "joepegs.com")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId2).WillReturnRows(mockedAssetRow2)

	repo := TransactionRepo{
		TxLeg: repository.NewTxLeg((sqlxDB)),
		User:  repository.NewUser(sqlxDB),
		Asset: repository.NewAsset(sqlxDB),
	}

	u21Transaction := NewTransaction(repo)

	pass, err := u21Transaction.Evaluate(transaction)
	assert.NoError(t, err)
	assert.True(t, pass)
}

// Entity makes a credit card purchase over $1,500
func TestEvaluateTransactionAbnormalAmounts(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	transactionId := uuid.NewString()
	OriginTxLegID := uuid.NewString()
	DestinationTxLegID := uuid.NewString()
	userId1 := uuid.NewString()
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	networkId := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          networkId,
		NetworkFee:         "100000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "2000000000",
		OriginTxLegID:      OriginTxLegID,
		ReceiptTxLegID:     sql.NullString{String: uuid.NewString()},
		ResponseTxLegID:    sql.NullString{String: uuid.NewString()},
		DestinationTxLegID: DestinationTxLegID,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(OriginTxLegID, time.Now(), "2000000000", "2000000000", assetId1, userId1, instrumentId1)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(OriginTxLegID).WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(DestinationTxLegID, time.Now(), "1", "2000000000", assetId2, userId1, instrumentId2)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(DestinationTxLegID).WillReturnRows(mockedTxLegRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId1, "USD", "fiat USD", 6, false, networkId, "self")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId1).WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId2, "Noose The Goose", "Noose the Goose NFT", 0, true, networkId, "joepegs.com")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId2).WillReturnRows(mockedAssetRow2)

	repo := TransactionRepo{
		TxLeg: repository.NewTxLeg((sqlxDB)),
		User:  repository.NewUser(sqlxDB),
		Asset: repository.NewAsset(sqlxDB),
	}

	u21Transaction := NewTransaction(repo)

	pass, err := u21Transaction.Evaluate(transaction)
	assert.NoError(t, err)
	assert.False(t, pass)
}

// User links more than 5 cards to their account in a 1 hour span
func TestEvaluateTransactionManyLinkedCards(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	userId1, u21UserId, err := createMockUser(mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21UserId)), 0)

	transactionId := uuid.NewString()
	OriginTxLegID := uuid.NewString()
	DestinationTxLegID := uuid.NewString()
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	networkId := uuid.NewString()
	var instrumentId1 string //:= uuid.NewString()
	instrumentId2 := uuid.NewString()

	// create 6 instruments
	for i := 0; i <= 5; i++ {
		var u21InstrumentId string
		instrument, u21InstrumentId, err := createMockInstrumentForUser(userId1, mock, sqlxDB)
		instrumentId1 = instrument.ID
		assert.NoError(t, err)
		assert.Greater(t, len([]rune(u21InstrumentId)), 0)

		// Log create instrument action w/ Unit21
		u21ActionRepo := ActionRepo{
			User:     repository.NewUser(sqlxDB),
			Device:   repository.NewDevice(sqlxDB),
			Location: repository.NewLocation(sqlxDB),
		}

		u21Action := NewAction(u21ActionRepo)
		_, err = u21Action.Create(instrument, "CreditCard", "Creation", u21InstrumentId, "Creation")
		if err != nil {
			fmt.Printf("Error creating a new instrument action in Unit21")
			return
		}
	}

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          networkId,
		NetworkFee:         "1000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "1000000",
		OriginTxLegID:      OriginTxLegID,
		ReceiptTxLegID:     sql.NullString{},
		ResponseTxLegID:    sql.NullString{},
		DestinationTxLegID: DestinationTxLegID,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(OriginTxLegID, time.Now(), "4000000", "4000000", assetId1, userId1, instrumentId1)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(OriginTxLegID).WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(DestinationTxLegID, time.Now(), "1", "4000000", assetId2, userId1, instrumentId2)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(DestinationTxLegID).WillReturnRows(mockedTxLegRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId1, "USD", "fiat USD", 6, false, networkId, "self")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId1).WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId2, "Noose The Goose", "Noose the Goose NFT", 0, true, networkId, "joepegs.com")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId2).WillReturnRows(mockedAssetRow2)

	repo := TransactionRepo{
		TxLeg: repository.NewTxLeg((sqlxDB)),
		User:  repository.NewUser(sqlxDB),
		Asset: repository.NewAsset(sqlxDB),
	}

	u21Transaction := NewTransaction(repo)

	pass, err := u21Transaction.Evaluate(transaction)
	assert.NoError(t, err)
	assert.False(t, pass)
}

// 10 or more FAILED transactions in a 1 hour span
func TestEvaluateTransactionHighFailedTransactionAmount(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	userId1, u21UserId, err := createMockUser(mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21UserId)), 0)

	transactionId := uuid.NewString()
	OriginTxLegID := uuid.NewString()
	DestinationTxLegID := uuid.NewString()
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	networkId := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          networkId,
		NetworkFee:         "100000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "2000000000",
		OriginTxLegID:      OriginTxLegID,
		ReceiptTxLegID:     sql.NullString{String: uuid.NewString()},
		ResponseTxLegID:    sql.NullString{String: uuid.NewString()},
		DestinationTxLegID: DestinationTxLegID,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(OriginTxLegID, time.Now(), "1000000", "1000000", assetId1, userId1, instrumentId1)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(OriginTxLegID).WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(DestinationTxLegID, time.Now(), "1", "10000000", assetId2, userId1, instrumentId2)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(DestinationTxLegID).WillReturnRows(mockedTxLegRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId1, "USD", "fiat USD", 6, false, networkId, "self")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId1).WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId2, "Noose The Goose", "Noose the Goose NFT", 0, true, networkId, "joepegs.com")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId2).WillReturnRows(mockedAssetRow2)

	repo := TransactionRepo{
		TxLeg: repository.NewTxLeg((sqlxDB)),
		User:  repository.NewUser(sqlxDB),
		Asset: repository.NewAsset(sqlxDB),
	}

	u21Transaction := NewTransaction(repo)

	pass, err := u21Transaction.Evaluate(transaction)
	assert.Error(t, err)
	assert.False(t, pass)
}

// User onboarded in the last 48 hours and has
// transacted more than 7.5K in the last 90 minutes
func TestEvaluateTransactionNewUserHighSpend(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	transactionId := uuid.NewString()
	OriginTxLegID := uuid.NewString()
	DestinationTxLegID := uuid.NewString()
	userId1 := uuid.NewString()
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	networkId := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          networkId,
		NetworkFee:         "100000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "2000000000",
		OriginTxLegID:      OriginTxLegID,
		ReceiptTxLegID:     sql.NullString{String: uuid.NewString()},
		ResponseTxLegID:    sql.NullString{String: uuid.NewString()},
		DestinationTxLegID: DestinationTxLegID,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(OriginTxLegID, time.Now(), "1000000", "1000000", assetId1, userId1, instrumentId1)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(OriginTxLegID).WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(DestinationTxLegID, time.Now(), "1", "10000000", assetId2, userId1, instrumentId2)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deactivated_at IS NULL").WithArgs(DestinationTxLegID).WillReturnRows(mockedTxLegRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId1, "USD", "fiat USD", 6, false, networkId, "self")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId1).WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId2, "Noose The Goose", "Noose the Goose NFT", 0, true, networkId, "joepegs.com")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deactivated_at IS NULL").WithArgs(assetId2).WillReturnRows(mockedAssetRow2)

	repo := TransactionRepo{
		TxLeg: repository.NewTxLeg((sqlxDB)),
		User:  repository.NewUser(sqlxDB),
		Asset: repository.NewAsset(sqlxDB),
	}

	u21Transaction := NewTransaction(repo)

	pass, err := u21Transaction.Evaluate(transaction)
	assert.Error(t, err)
	assert.False(t, pass)
}
