package unit21

import (
	"context"
	"testing"
	"time"

	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// This transaction should pass
func TestEvaluateTransactionPass(t *testing.T) {
	config.LoadEnv("../../../.env")
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
	pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
	assert.NoError(t, err)
	assert.True(t, pass)
}

// Entity makes a credit card purchase over $1,500
func TestEvaluateTransactionAbnormalAmounts(t *testing.T) {
	config.LoadEnv("../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId := uuid.NewString()
	transaction := createMockTransactionForUser(userId, "2000000000", sqlxDB)
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
	pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
	assert.NoError(t, err)
	assert.False(t, pass)
}

// User links more than 5 cards to their account in a 1 hour span
// Not currently functioning due to lag in Unit21 data ingestion
func TestEvaluateTransactionManyLinkedCards(t *testing.T) {
	config.LoadEnv("../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId, u21UserId, err := createMockUser(mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21UserId)), 0)

	// create 6 instruments
	for i := 0; i <= 5; i++ {
		var u21InstrumentId string
		instrument, u21InstrumentId, err := createMockInstrumentForUser(userId, mock, sqlxDB)
		assert.NoError(t, err)
		assert.Greater(t, len([]rune(u21InstrumentId)), 0)

		u21Action := NewAction()
		_, err = u21Action.Create(instrument, "Creation", u21InstrumentId, "Creation")
		if err != nil {
			t.Log("Error creating a new instrument action in Unit21")
			return
		}
	}

	transaction := createMockTransactionForUser(userId, "1000000", sqlxDB)
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
	time.Sleep(10 * time.Second)
	pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
	assert.NoError(t, err)
	assert.False(t, pass)
}

// 10 or more FAILED transactions in a 1 hour span
func TestEvaluateTransactionHighFailedTransactionAmount(t *testing.T) {
	config.LoadEnv("../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId, u21UserId, err := createMockUser(mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21UserId)), 0)
	// create, evaluate, and execute 10 failed transactions
	for i := 0; i < 10; i++ {
		transaction := createMockTransactionForUser(userId, "2000000000", sqlxDB)
		assetId1 := uuid.NewString()
		assetId2 := uuid.NewString()
		instrumentId1 := uuid.NewString()
		instrumentId2 := uuid.NewString()
		mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
		pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
		assert.NoError(t, err)
		assert.False(t, pass)
		transaction.Status = "Failed"
		mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
		u21TransactionId, err := executeMockTransactionForUser(ctx, transaction, sqlxDB)
		assert.NoError(t, err)
		assert.Greater(t, len([]rune(u21TransactionId)), 0)
	}

	// create and evaluate a transaction that would otherwise pass
	transaction := createMockTransactionForUser(userId, "1000000", sqlxDB)
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
	time.Sleep(10 * time.Second)
	pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
	assert.NoError(t, err)
	assert.False(t, pass)
}

// User onboarded in the last 48 hours and has
// transacted more than 7.5K in the last 90 minutes
func TestEvaluateTransactionNewUserHighSpend(t *testing.T) {
	config.LoadEnv("../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId, u21UserId, err := createMockUser(mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21UserId)), 0)
	// create, evaluate, and execute 6 successful transactions for $1500
	for i := 0; i < 6; i++ {
		transaction := createMockTransactionForUser(userId, "1500000000", sqlxDB)
		assetId1 := uuid.NewString()
		assetId2 := uuid.NewString()
		instrumentId1 := uuid.NewString()
		instrumentId2 := uuid.NewString()
		mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
		pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
		assert.NoError(t, err)
		assert.True(t, pass)
		mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
		u21TransactionId, err := executeMockTransactionForUser(ctx, transaction, sqlxDB)
		assert.NoError(t, err)
		assert.Greater(t, len([]rune(u21TransactionId)), 0)
	}

	// create and evaluate a transaction that would otherwise pass
	transaction := createMockTransactionForUser(userId, "1000000", sqlxDB)
	assetId1 := uuid.NewString()
	assetId2 := uuid.NewString()
	instrumentId1 := uuid.NewString()
	instrumentId2 := uuid.NewString()
	mockTransactionRows(mock, transaction, userId, assetId1, assetId2, instrumentId1, instrumentId2)
	time.Sleep(10 * time.Second)
	pass, err := evaluateMockTransaction(ctx, transaction, sqlxDB)
	assert.NoError(t, err)
	assert.False(t, pass)
}

func evaluateMockTransaction(ctx context.Context, transaction model.Transaction, sqlxDB *sqlx.DB) (pass bool, err error) {
	repos := TransactionRepos{
		TxLeg:  repository.NewTxLeg((sqlxDB)),
		User:   repository.NewUser(sqlxDB),
		Asset:  repository.NewAsset(sqlxDB),
		Device: repository.NewDevice(sqlxDB),
	}

	u21Transaction := NewTransaction(repos)

	pass, err = u21Transaction.Evaluate(ctx, transaction)

	return
}
