package unit21

import (
	"database/sql"
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

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	_, u21EntityId, err := createMockUser(mock, sqlxDB)

	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// TODO: mock call to client once it's manually tested
}

func TestUpdateEntity(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	// create entity to modify
	entityId, u21EntityId, err := createMockUser(mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

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
		AddRow(uuid.NewString(), entityId, "email", "verified", "test@gmail.com").AddRow(uuid.NewString(), entityId, "phone", "verified", "+12345678910")
	mock.ExpectQuery("SELECT * FROM contact WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedContactRow)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(uuid.NewString(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, entityId)
	mock.ExpectQuery("SELECT * FROM device WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedDeviceRow)

	mockedUserPlatformRow := sqlmock.NewRows([]string{"user_id", "platform_id"}).
		AddRow(entityId, uuid.NewString())
	mock.ExpectQuery("SELECT * FROM user_to_platform WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedUserPlatformRow)

	repos := EntityRepos{
		Device:         repository.NewDevice(sqlxDB),
		Contact:        repository.NewContact(sqlxDB),
		UserToPlatform: repository.NewUserToPlatform(sqlxDB),
	}

	u21Entity := NewEntity(repos)

	// update in u21
	u21EntityId, err = u21Entity.Update(user)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// TODO: mock call to client once it's manually tested
}

func TestAddInstruments(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	// create entity to modify
	entityId, u21EntityId, err := createMockUser(mock, sqlxDB)

	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

	var instrumentIds []string

	// mock new instrumentIds
	for i := 1; i <= 10; i++ {
		instrumentIds = append(instrumentIds, uuid.NewString())
	}

	repos := EntityRepos{
		Device:         repository.NewDevice(sqlxDB),
		Contact:        repository.NewContact(sqlxDB),
		UserToPlatform: repository.NewUserToPlatform(sqlxDB),
	}

	u21Entity := NewEntity(repos)
	err = u21Entity.AddInstruments(entityId, instrumentIds)
	assert.NoError(t, err)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// TODO: mock call to client once it's manually tested
}

func createMockUser(mock sqlmock.Sqlmock, sqlxDB *sqlx.DB) (entityId string, unit21Id string, err error) {
	entityId = uuid.NewString()
	user := model.User{
		ID:            entityId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeactivatedAt: nil,
		Type:          "User",
		Status:        "Onboarded",
		Tags:          model.StringMap{"platform": "Activision Blizzard"},
		FirstName:     "Test",
		MiddleName:    "A",
		LastName:      "User",
	}

	mockedContactRow := sqlmock.NewRows([]string{"id", "user_id", "type", "status", "data"}).
		AddRow(uuid.NewString(), entityId, "email", "verified", "test@gmail.com").AddRow(uuid.NewString(), entityId, "phone", "verified", "+12345678910")
	mock.ExpectQuery("SELECT * FROM contact WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedContactRow)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(uuid.NewString(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, entityId)
	mock.ExpectQuery("SELECT * FROM device WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedDeviceRow)

	mockedUserPlatformRow := sqlmock.NewRows([]string{"user_id", "platform_id"}).
		AddRow(entityId, uuid.NewString())
	mock.ExpectQuery("SELECT * FROM user_to_platform WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(entityId, 100, 0).WillReturnRows(mockedUserPlatformRow)

	repos := EntityRepos{
		Device:         repository.NewDevice(sqlxDB),
		Contact:        repository.NewContact(sqlxDB),
		UserToPlatform: repository.NewUserToPlatform(sqlxDB),
	}

	u21Entity := NewEntity(repos)

	u21EntityId, err := u21Entity.Create(user)

	return entityId, u21EntityId, err
}

func createMockInstrumentForUser(userId string, mock sqlmock.Sqlmock, sqlxDB *sqlx.DB) (instrument model.Instrument, unit21Id string, err error) {
	instrumentId := uuid.NewString()
	locationId := uuid.NewString()

	instrument = model.Instrument{
		ID:            instrumentId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeactivatedAt: nil,
		Type:          "Credit Card",
		Status:        "Verified",
		Tags:          nil,
		Network:       "Visa",
		PublicKey:     "",
		Last4:         "1234",
		UserID:        userId,
		LocationID:    sql.NullString{String: locationId},
	}

	mockedUserRow1 := sqlmock.NewRows([]string{"id", "type", "status", "tags", "first_name", "middle_name", "last_name"}).
		AddRow(userId, "User", "Onboarded", `{"kyc_level": "1", "platform": "mortal kombat"}`, "Daemon", "", "Targaryan")
	mock.ExpectQuery("SELECT * FROM string_user WHERE id = $1 AND deactivated_at IS NULL").WithArgs(userId).WillReturnRows(mockedUserRow1)

	mockedUserRow2 := sqlmock.NewRows([]string{"id", "type", "status", "tags", "first_name", "middle_name", "last_name"}).
		AddRow(userId, "User", "Onboarded", `{"kyc_level": "1", "platform": "mortal kombat"}`, "Daemon", "", "Targaryan")
	mock.ExpectQuery("SELECT * FROM string_user WHERE id = $1 AND deactivated_at IS NULL").WithArgs(userId).WillReturnRows(mockedUserRow2)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(uuid.NewString(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, userId)
	mock.ExpectQuery("SELECT * FROM device WHERE user_id = $1 LIMIT $2 OFFSET $3").WithArgs(userId, 100, 0).WillReturnRows(mockedDeviceRow)

	mockedLocationRow := sqlmock.NewRows([]string{"id", "type", "status", "building_number", "unit_number", "street_name", "city", "state", "postal_code", "country"}).
		AddRow(locationId, "Home", "Verified", "20181", "411", "Lark Avenue", "Somerville", "MA", "01443", "USA")
	mock.ExpectQuery("SELECT * FROM location WHERE id = $1 AND deactivated_at IS NULL").WithArgs(locationId).WillReturnRows(mockedLocationRow)

	repo := InstrumentRepo{
		User:     repository.NewUser(sqlxDB),
		Device:   repository.NewDevice(sqlxDB),
		Location: repository.NewLocation(sqlxDB),
	}

	u21Instrument := NewInstrument(repo)

	u21InstrumentId, err := u21Instrument.Create(instrument)

	return instrument, u21InstrumentId, err
}
