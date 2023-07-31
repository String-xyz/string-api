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
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstrument(t *testing.T) {
	env.LoadEnv(&config.Var, "../../../.env")
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId := uuid.NewString()

	_, u21InstrumentId, err := createMockInstrumentForUser(userId, mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

}

func TestUpdateInstrument(t *testing.T) {
	env.LoadEnv(&config.Var, "../../../.env")
	ctx := context.Background()
	db, mock, sqlxDB, err := initializeTest(t)
	assert.NoError(t, err)
	defer db.Close()

	userId := uuid.NewString()

	instrument, u21InstrumentId, err := createMockInstrumentForUser(userId, mock, sqlxDB)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	locationId := uuid.NewString()

	instrument = model.Instrument{
		Id:         instrument.Id,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Type:       "Credit Card",
		Status:     "Verified",
		Tags:       nil,
		Network:    "Visa",
		PublicKey:  "",
		Last4:      "1235",
		UserId:     userId,
		LocationId: sql.NullString{String: locationId},
	}

	mockedUserRow1 := sqlmock.NewRows([]string{"id", "type", "status", "tags", "first_name", "middle_name", "last_name"}).
		AddRow(userId, "User", "Onboarded", `{"kyc_level": "1", "platform": "mortal kombat"}`, "Daemon", "", "Targaryan")
	mock.ExpectQuery("SELECT * FROM string_user WHERE id = $1 AND deleted_at IS NULL").WithArgs(userId).WillReturnRows(mockedUserRow1)

	mockedUserRow2 := sqlmock.NewRows([]string{"id", "type", "status", "tags", "first_name", "middle_name", "last_name"}).
		AddRow(userId, "User", "Onboarded", `{"kyc_level": "1", "platform": "mortal kombat"}`, "Daemon", "", "Targaryan")
	mock.ExpectQuery("SELECT * FROM string_user WHERE id = $1 AND deleted_at IS NULL").WithArgs(userId).WillReturnRows(mockedUserRow2)

	mockedDeviceRow := sqlmock.NewRows([]string{"id", "type", "description", "fingerprint", "ip_addresses", "user_id"}).
		AddRow(uuid.NewString(), "Mobile", "iPhone 11S", uuid.NewString(), pq.StringArray{"187.25.24.128"}, userId)
	mock.ExpectQuery("SELECT * FROM device WHERE user_id = $1 AND deleted_at IS NULL LIMIT $2 OFFSET $3").WithArgs(userId, 100, 0).WillReturnRows(mockedDeviceRow)

	mockedLocationRow := sqlmock.NewRows([]string{"id", "type", "status", "building_number", "unit_number", "street_name", "city", "state", "postal_code", "country"}).
		AddRow(locationId, "Home", "Verified", "20181", "411", "Lark Avenue", "Somerville", "MA", "01443", "USA")
	mock.ExpectQuery("SELECT * FROM location WHERE id = $1 AND deleted_at IS NULL").WithArgs(locationId).WillReturnRows(mockedLocationRow)

	repos := InstrumentRepos{
		User:     repository.NewUser(sqlxDB),
		Device:   repository.NewDevice(sqlxDB),
		Location: repository.NewLocation(sqlxDB),
	}

	action := NewAction()

	u21Instrument := NewInstrument(repos, action)

	u21InstrumentId, err = u21Instrument.Update(ctx, instrument)
	assert.NoError(t, err)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new instrument added
	// TODO: mock call to client once it's manually tested
}
