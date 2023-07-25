package unit21

import (
	"context"
	"database/sql"
	"testing"
	"time"

	env "github.com/String-xyz/go-lib/v2/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	env.LoadEnv(&config.Var, "../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId := uuid.NewString()
	transaction := createMockTransactionForUser(userId, "1000000", sqlxDB)
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
	u21TransactionId, err := executeMockTransactionForUser(ctx, transaction, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21TransactionId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new transaction added
	// TODO: mock call to client once it's manually tested
}

func TestUpdateTransaction(t *testing.T) {
	env.LoadEnv(&config.Var, "../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId := uuid.NewString()
	transaction := createMockTransactionForUser(userId, "1000000", sqlxDB)
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
	u21TransactionId, err := executeMockTransactionForUser(ctx, transaction, sqlxDB)
	assert.NoError(t, err)

	OriginTxLegId := uuid.NewString()
	DestinationTxLegId := uuid.NewString()
	assetId1 = uuid.NewString()
	assetId2 = uuid.NewString()
	networkId := uuid.NewString()
	instrumentId1 = uuid.NewString()
	instrumentId2 = uuid.NewString()

	transaction = model.Transaction{
		Id:                 transaction.Id,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceId:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformId:         uuid.NewString(),
		TransactionHash:    "",
		NetworkId:          networkId,
		NetworkFee:         "100000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  "1000000000",
		OriginTxLegId:      OriginTxLegId,
		ReceiptTxLegId:     sql.NullString{String: uuid.NewString()},
		ResponseTxLegId:    sql.NullString{String: uuid.NewString()},
		DestinationTxLegId: DestinationTxLegId,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "2000000",
	}
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)

	repos := TransactionRepos{
		TxLeg:  repository.NewTxLeg((sqlxDB)),
		User:   repository.NewUser(sqlxDB),
		Asset:  repository.NewAsset(sqlxDB),
		Device: repository.NewDevice(sqlxDB),
	}

	u21Transaction := NewTransaction(repos)

	u21TransactionId, err = u21Transaction.Update(ctx, transaction)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21TransactionId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new transaction added
	// TODO: mock call to client once it's manually tested
}

func executeMockTransactionForUser(ctx context.Context, transaction model.Transaction, sqlxDB *sqlx.DB) (unit21Id string, err error) {
	repos := TransactionRepos{
		TxLeg:  repository.NewTxLeg(sqlxDB),
		User:   repository.NewUser(sqlxDB),
		Asset:  repository.NewAsset(sqlxDB),
		Device: repository.NewDevice(sqlxDB),
	}

	u21Transaction := NewTransaction(repos)

	unit21Id, err = u21Transaction.Create(ctx, transaction)

	return
}

func createMockTransactionForUser(userId string, amount string, sqlxDB *sqlx.DB) (transaction model.Transaction) {
	transactionId := uuid.NewString()
	OriginTxLegId := uuid.NewString()
	DestinationTxLegId := uuid.NewString()
	networkId := uuid.NewString()

	transaction = model.Transaction{
		Id:                 transactionId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Type:               "fiat-to-crypto",
		Status:             "Completed",
		Tags:               map[string]string{},
		DeviceId:           uuid.NewString(),
		IPAddress:          "187.25.24.128",
		PlatformId:         uuid.NewString(),
		TransactionHash:    "",
		NetworkId:          networkId,
		NetworkFee:         "100000000",
		ContractParams:     pq.StringArray{},
		ContractFunc:       "mintTo()",
		TransactionAmount:  amount,
		OriginTxLegId:      OriginTxLegId,
		ReceiptTxLegId:     sql.NullString{String: uuid.NewString()},
		ResponseTxLegId:    sql.NullString{String: uuid.NewString()},
		DestinationTxLegId: DestinationTxLegId,
		ProcessingFee:      "1000000",
		ProcessingFeeAsset: uuid.NewString(),
		StringFee:          "1000000",
	}

	return
}

func mockTransactionRows(mock sqlmock.Sqlmock, transaction model.Transaction, userId string, assetId1 string, assetId2 string, instrumentId1 string, instrumentId2 string) {
	mockedTxLegRow1 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(transaction.OriginTxLegId, time.Now(), transaction.TransactionAmount, transaction.TransactionAmount, assetId1, userId, instrumentId1)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deleted_at IS NULL").WithArgs(transaction.OriginTxLegId).WillReturnRows(mockedTxLegRow1)

	mockedTxLegRow2 := sqlmock.NewRows([]string{"id", "timestamp", "amount", "value", "asset_id", "user_id", "instrument_id"}).
		AddRow(transaction.DestinationTxLegId, time.Now(), "1", transaction.TransactionAmount, assetId2, userId, instrumentId2)
	mock.ExpectQuery("SELECT * FROM tx_leg WHERE id = $1 AND deleted_at IS NULL").WithArgs(transaction.DestinationTxLegId).WillReturnRows(mockedTxLegRow2)

	mockedAssetRow1 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId1, "USD", "fiat USD", 6, false, transaction.NetworkId, "self")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deleted_at IS NULL").WithArgs(assetId1).WillReturnRows(mockedAssetRow1)

	mockedAssetRow2 := sqlmock.NewRows([]string{"id", "name", "description", "decimals", "is_crypto", "network_id", "value_oracle"}).
		AddRow(assetId2, "Noose The Goose", "Noose the Goose NFT", 0, true, transaction.NetworkId, "joepegs.com")
	mock.ExpectQuery("SELECT * FROM asset WHERE id = $1 AND deleted_at IS NULL").WithArgs(assetId2).WillReturnRows(mockedAssetRow2)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(transaction.DeviceId, "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, userId)
	mock.ExpectQuery("SELECT * FROM device WHERE id = $1 AND deleted_at IS NULL").WithArgs(transaction.DeviceId).WillReturnRows(mockedDeviceRow)

}
