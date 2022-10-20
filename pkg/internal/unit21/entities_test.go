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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntity(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	entityId := uuid.NewString()
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	user := model.User{
		ID:            entityId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeactivatedAt: nil,
		Type:          "User",
		Status:        "Onboarded",
		Tags:          nil,
		FirstName:     "Test",
		MiddleName:    "A",
		LastName:      "User",
	}

	// "SELECT * FROM %s WHERE user_id = $1 LIMIT $2 OFFSET $3", b.table), userID, limit, offset)
	mock.ExpectQuery("SELECT * FROM (.+) WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)")

	// Dependent on Device and Instrument Repos being created
	userRepo := repository.NewUser(sqlxDB)
	// deviceRepo := repository.NewDevice(sqlxDB)
	contactRepo := repository.NewUserContact(sqlxDB)
	// instrumentRepo := repository.NewInstrument(sqlxDB)

	// u21Entity := newEntity(userRepo, deviceRepo, contactRepo, instrumentRepo)
	u21Entity := newEntity(userRepo, contactRepo)

	u21EntityId, err := u21Entity.Create(user)
	assert.NoError(t, err)
	log.Printf("u21EntityId: %s", u21EntityId)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// todo: mock call to client once it's manually tested
}

func TestAddInstruments(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	entityId := "44142758-f015-4f79-a004-e554b0641480" //previous created test user
	var instrumentIds []string
	db, _, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	// mock new instrumentIds
	for i := 1; i <= 10; i++ {
		instrumentIds = append(instrumentIds, uuid.NewString())
	}

	// Dependent on Device and Instrument Repos being created
	userRepo := repository.NewUser(sqlxDB)
	// deviceRepo := repository.NewDevice(sqlxDB)
	contactRepo := repository.NewUserContact(sqlxDB)
	// instrumentRepo := repository.NewInstrument(sqlxDB)

	// u21Entity := newEntity(userRepo, deviceRepo, contactRepo, instrumentRepo)
	u21Entity := newEntity(userRepo, contactRepo)

	err = u21Entity.AddInstruments(entityId, instrumentIds)
	assert.NoError(t, err)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// todo: mock call to client once it's manually tested
}
