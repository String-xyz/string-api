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

func TestCreateInstrument(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	userId := uuid.NewString()

	_, u21InstrumentId, err := createMockInstrumentForUser(userId, mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

}

func TestUpdateInstrument(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	userId := uuid.NewString()

	instrument, u21InstrumentId, err := createMockInstrumentForUser(userId, mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	locationId := uuid.NewString()

	instrument = model.Instrument{
		ID:            instrument.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeactivatedAt: nil,
		Type:          "Credit Card",
		Status:        "Verified",
		Tags:          nil,
		Network:       "Visa",
		PublicKey:     "",
		Last4:         "1235",
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

	u21InstrumentId, err = u21Instrument.Update(instrument)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new instrument added
	// TODO: mock call to client once it's manually tested
}
