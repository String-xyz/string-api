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
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntity(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	entityId := uuid.NewString()
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

	mockedContactRow := sqlmock.NewRows([]string{"id", "user_id", "type", "status", "data"}).
		AddRow(uuid.NewString(), entityId, "email", "verified", "test@gmail.com")
	mock.ExpectQuery("SELECT * FROM contact WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedContactRow)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(uuid.NewString(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, entityId)
	mock.ExpectQuery("SELECT * FROM device WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedDeviceRow)

	mockedUserPlatformRow := sqlmock.NewRows([]string{"user_id", "platform_id"}).
		AddRow(entityId, uuid.NewString())
	mock.ExpectQuery("SELECT * FROM user_platform WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedUserPlatformRow)

	repos := EntityRepos{
		device:       repository.NewDevice(sqlxDB),
		contact:      repository.NewContact(sqlxDB),
		userPlatform: repository.NewUserPlatform(sqlxDB),
	}

	u21Entity := NewEntity(repos)

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

	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	// choose an older entity so we can update it
	entityId := "d8451ddd-6116-4b62-9072-0e8b63a843f1"
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

	mockedContactRow := sqlmock.NewRows([]string{"id", "user_id", "type", "status", "data"}).
		AddRow(uuid.NewString(), entityId, "email", "verified", "test@gmail.com")
	mock.ExpectQuery("SELECT * FROM contact WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedContactRow)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(uuid.NewString(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, entityId)
	mock.ExpectQuery("SELECT * FROM device WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedDeviceRow)

	mockedUserPlatformRow := sqlmock.NewRows([]string{"user_id", "platform_id"}).
		AddRow(entityId, uuid.NewString())
	mock.ExpectQuery("SELECT * FROM user_platform WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedUserPlatformRow)

	repos := EntityRepos{
		device:       repository.NewDevice(sqlxDB),
		contact:      repository.NewContact(sqlxDB),
		userPlatform: repository.NewUserPlatform(sqlxDB),
	}

	u21Entity := NewEntity(repos)

	// update in u21
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

	db, _, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	entityId := "44142758-f015-4f79-a004-e554b0641480" //previous created test user
	var instrumentIds []string

	// mock new instrumentIds
	for i := 1; i <= 10; i++ {
		instrumentIds = append(instrumentIds, uuid.NewString())
	}

	repos := EntityRepos{
		device:       repository.NewDevice(sqlxDB),
		contact:      repository.NewContact(sqlxDB),
		userPlatform: repository.NewUserPlatform(sqlxDB),
	}

	u21Entity := NewEntity(repos)

	err = u21Entity.AddInstruments(entityId, instrumentIds)
	assert.NoError(t, err)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// todo: mock call to client once it's manually tested
}
