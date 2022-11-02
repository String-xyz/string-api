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

func TestCreateInstrument(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	instrumentId := uuid.NewString()
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	instrument := model.Instrument{
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
		UserID:        uuid.NewString(),
		LocationID:    uuid.NewString(),
	}

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "last_used_at", "validated_at", "type", "description", "fingerprint", "ip_addresses", "user_id"}).AddRow(uuid.NewString(), time.Now(), time.Now(), time.Now(), time.Now(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"192.0.1.1"}, uuid.NewString())
	mock.ExpectQuery(`SELECT \* FROM device WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedDeviceRow)

	mockedLocationRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "status", "tags", "building_numeber", "unit_number", "street_name", "city", "state", "postal_code", "country"}).AddRow(uuid.NewString(), time.Now(), time.Now(), "Home", "Verified", nil, "20181", "411", "Lark Avenue", "Somerville", "MA", "01443", "USA")
	mock.ExpectQuery(`SELECT \* FROM location WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedLocationRow)

	repo := InstrumentRepo{
		user:     repository.NewUser(sqlxDB),
		device:   repository.NewDevice(sqlxDB),
		location: repository.NewLocation(sqlxDB),
	}

	u21Instrument := NewInstrument(repo)

	u21InstrumentId, err := u21Instrument.Create(instrument)
	assert.NoError(t, err)
	log.Printf("u21InstrumentId: %s", u21InstrumentId)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new instrument added
	// todo: mock call to client once it's manually tested
}

func TestUpdateInstrument(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	instrumentId := "bb73ef2d-0e62-4381-8a37-6a32c11fb226"
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	instrument := model.Instrument{
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
		UserID:        uuid.NewString(),
		LocationID:    uuid.NewString(),
	}

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "last_used_at", "validated_at", "type", "description", "fingerprint", "ip_addresses", "user_id"}).AddRow(uuid.NewString(), time.Now(), time.Now(), time.Now(), time.Now(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"192.0.1.1"}, uuid.NewString())
	mock.ExpectQuery(`SELECT \* FROM device WHERE user_id = (.+) LIMIT (.+) OFFSET (.+)`).WithArgs().WillReturnRows(mockedDeviceRow)

	mockedLocationRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "status", "tags", "building_numeber", "unit_number", "street_name", "city", "state", "postal_code", "country"}).AddRow(uuid.NewString(), time.Now(), time.Now(), "Home", "Verified", nil, "20181", "411", "Lark Avenue", "Somerville", "MA", "01443", "USA")
	mock.ExpectQuery(`SELECT \* FROM location WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedLocationRow)

	repo := InstrumentRepo{
		user:     repository.NewUser(sqlxDB),
		device:   repository.NewDevice(sqlxDB),
		location: repository.NewLocation(sqlxDB),
	}

	u21Instrument := NewInstrument(repo)

	u21InstrumentId, err := u21Instrument.Update(instrument)
	assert.NoError(t, err)
	log.Printf("u21InstrumentId: %s", u21InstrumentId)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new instrument added
	// todo: mock call to client once it's manually tested
}
