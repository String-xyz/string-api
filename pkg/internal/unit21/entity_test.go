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

	mockedContactRow := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at", "last_authenticated_at", "deactivated_at", "type", "status", "data"})
	//.AddRow(1, time.Now(), time.Now(), "1")
	mock.ExpectQuery(`SELECT \* FROM contact WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedContactRow)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "last_used_at", "validated_at", "deactivated_at", "type", "description", "fingerprint", "ip_addresses", "user_id"})
	mock.ExpectQuery(`SELECT \* FROM device WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedDeviceRow)

	mockedUserPlatformRow := sqlmock.NewRows([]string{"user_id", "platform_id"})
	mock.ExpectQuery(`SELECT \* FROM user_platform WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedUserPlatformRow)

	// Dependent on Device and Instrument Repos being created
	userRepo := repository.NewUser(sqlxDB)
	deviceRepo := repository.NewDevice(sqlxDB)
	contactRepo := repository.NewContact(sqlxDB)
	userPlatformRepo := repository.NewUserPlatform(sqlxDB)

	u21Entity := newEntity(userRepo, deviceRepo, contactRepo, userPlatformRepo)

	u21EntityId, err := u21Entity.Create(user)
	assert.NoError(t, err)
	log.Printf("u21EntityId: %s", u21EntityId)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// todo: mock call to client once it's manually tested
}

func TestUpdateEntity(t *testing.T) {
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

	mockedContactRow := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at", "last_authenticated_at", "deactivated_at", "type", "status", "data"})
	//.AddRow(1, time.Now(), time.Now(), "1")
	mock.ExpectQuery(`SELECT \* FROM contact WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedContactRow)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "last_used_at", "validated_at", "deactivated_at", "type", "description", "fingerprint", "ip_addresses", "user_id"})
	mock.ExpectQuery(`SELECT \* FROM device WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedDeviceRow)

	mockedUserPlatformRow := sqlmock.NewRows([]string{"user_id", "platform_id"})
	mock.ExpectQuery(`SELECT \* FROM user_platform WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedUserPlatformRow)

	// Dependent on Device and Instrument Repos being created
	userRepo := repository.NewUser(sqlxDB)
	deviceRepo := repository.NewDevice(sqlxDB)
	contactRepo := repository.NewContact(sqlxDB)
	userPlatformRepo := repository.NewUserPlatform(sqlxDB)

	u21Entity := newEntity(userRepo, deviceRepo, contactRepo, userPlatformRepo)

	u21EntityId, err := u21Entity.Update(user)
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
	deviceRepo := repository.NewDevice(sqlxDB)
	contactRepo := repository.NewContact(sqlxDB)
	userPlatformRepo := repository.NewUserPlatform(sqlxDB)

	u21Entity := newEntity(userRepo, deviceRepo, contactRepo, userPlatformRepo)

	err = u21Entity.AddInstruments(entityId, instrumentIds)
	assert.NoError(t, err)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// todo: mock call to client once it's manually tested
}
