package unit21

import (
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	transactionId := uuid.NewString()
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "192.0.1.1",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          uuid.NewString(),
		NetworkFee:         "100000000",
		ContractParams:     types.JSONText{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "1000000000",
		OriginTxLegID:      uuid.NewString(),
		ReceiptTxLegID:     uuid.NewString(),
		ResponseTxLegID:    uuid.NewString(),
		DestinationTxLegID: uuid.NewString(),
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).AddRow(uuid.NewString(), time.Now(), time.Now(), time.Now(), "1000000", "1000000", uuid.NewString(), uuid.NewString(), uuid.NewString())
	mock.ExpectQuery(`SELECT \* FROM tx_leg WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).AddRow(uuid.NewString(), time.Now(), time.Now(), time.Now(), "10000000", "10000000", uuid.NewString(), uuid.NewString(), uuid.NewString())
	mock.ExpectQuery(`SELECT \* FROM tx_leg WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedTxLegRow2)

	mockedUserRow1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deactivated_at", "type", "status", "tags", "first_name", "middle_name", "last_name"}).AddRow(uuid.NewString(), time.Now(), time.Now(), nil, "User", "Onboarded", `{"kyc_level": "1", "platform": "mortal kombat"}`, "Daemon", "", "Targaryan")
	mock.ExpectQuery(`SELECT \* FROM user WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedUserRow1)

	mockedUserRow2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deactivated_at", "type", "status", "tags", "first_name", "middle_name", "last_name"}).AddRow(uuid.NewString(), time.Now(), time.Now(), nil, "User", "Onboarded", `{"kyc_level": "1", "platform": "space invaders"}`, "Toph", "", "Bei Fong")
	mock.ExpectQuery(`SELECT \* FROM user WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedUserRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).AddRow(uuid.NewString(), time.Now(), time.Now(), "USD", "fiat USD", 6, false, uuid.NewString(), "self")
	mock.ExpectQuery(`SELECT \* FROM asset WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).AddRow(uuid.NewString(), time.Now(), time.Now(), "Noose The Goose", "Noose the Goose NFT", 1, true, uuid.NewString(), "joepegs.com")
	mock.ExpectQuery(`SELECT \* FROM asset WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedAssetRow2)

	transactionRepo := repository.NewTransaction(sqlxDB)
	txLegRepo := repository.NewTxLeg((sqlxDB))
	userRepo := repository.NewUser(sqlxDB)
	assetRepo := repository.NewAsset(sqlxDB)

	u21Transaction := NewTransaction(transactionRepo, txLegRepo, userRepo, assetRepo)

	u21TransactionId, err := u21Transaction.Create(transaction)
	assert.NoError(t, err)
	log.Printf("u21TransactionId: %s", u21TransactionId)
	assert.Greater(t, len([]rune(u21TransactionId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new transaction added
	// todo: mock call to client once it's manually tested
}

func TestUpdateTransaction(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	transactionId := "866c32ae-9fda-409c-a0be-28b830022e93"
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	transaction := model.Transaction{
		ID:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceID:           uuid.NewString(),
		IPAddress:          "192.206.151.131",
		PlatformID:         uuid.NewString(),
		TransactionHash:    "",
		NetworkID:          uuid.NewString(),
		NetworkFee:         "100000000",
		ContractParams:     types.JSONText{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "1000000000",
		OriginTxLegID:      uuid.NewString(),
		ReceiptTxLegID:     uuid.NewString(),
		ResponseTxLegID:    uuid.NewString(),
		DestinationTxLegID: uuid.NewString(),
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "2000000",
	}

	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).AddRow(uuid.NewString(), time.Now(), time.Now(), time.Now(), "1000000", "1000000", uuid.NewString(), uuid.NewString(), uuid.NewString())
	mock.ExpectQuery(`SELECT \* FROM tx_leg WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).AddRow(uuid.NewString(), time.Now(), time.Now(), time.Now(), "10000000", "10000000", uuid.NewString(), uuid.NewString(), uuid.NewString())
	mock.ExpectQuery(`SELECT \* FROM tx_leg WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedTxLegRow2)

	mockedUserRow1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deactivated_at", "type", "status", "tags", "first_name", "middle_name", "last_name"}).AddRow(uuid.NewString(), time.Now(), time.Now(), nil, "User", "Onboarded", `{"kyc_level": "1", "platform": "mortal kombat"}`, "Daemon", "", "Targaryan")
	mock.ExpectQuery(`SELECT \* FROM user WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedUserRow1)

	mockedUserRow2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deactivated_at", "type", "status", "tags", "first_name", "middle_name", "last_name"}).AddRow(uuid.NewString(), time.Now(), time.Now(), nil, "User", "Onboarded", `{"kyc_level": "1", "platform": "space invaders"}`, "Toph", "", "Bei Fong")
	mock.ExpectQuery(`SELECT \* FROM user WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedUserRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).AddRow(uuid.NewString(), time.Now(), time.Now(), "USD", "fiat USD", 6, false, uuid.NewString(), "self")
	mock.ExpectQuery(`SELECT \* FROM asset WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).AddRow(uuid.NewString(), time.Now(), time.Now(), "Noose The Goose", "Noose the Goose NFT", 1, true, uuid.NewString(), "joepegs.com")
	mock.ExpectQuery(`SELECT \* FROM asset WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedAssetRow2)

	transactionRepo := repository.NewTransaction(sqlxDB)
	txLegRepo := repository.NewTxLeg((sqlxDB))
	userRepo := repository.NewUser(sqlxDB)
	assetRepo := repository.NewAsset(sqlxDB)

	u21Transaction := NewTransaction(transactionRepo, txLegRepo, userRepo, assetRepo)

	u21TransactionId, err := u21Transaction.Update(transaction)
	assert.NoError(t, err)
	log.Printf("u21TransactionId: %s", u21TransactionId)
	assert.Greater(t, len([]rune(u21TransactionId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new transaction added
	// todo: mock call to client once it's manually tested
}
