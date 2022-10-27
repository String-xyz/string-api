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

	mockedUserRow := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at", "deactivated_at", "type", "status", "tags", "first_name", "middle_name", "last_name"})
	//.AddRow(1, time.Now(), time.Now(), "1")
	mock.ExpectQuery(`SELECT string_user FROM %s WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedUserRow)

	mockedLocationRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "status", "tags", "building_numeber", "unit_number", "street_name", "city", "state", "postal_code", "country"})
	mock.ExpectQuery(`SELECT location FROM %s WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedLocationRow)

	instrumentRepo := repository.NewInstrument(sqlxDB)
	userRepo := repository.NewUser(sqlxDB)
	locationRepo := repository.NewLocation(sqlxDB)

	u21Instrument := NewInstrument(instrumentRepo, userRepo, locationRepo)

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

	mockedUserRow := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at", "deactivated_at", "type", "status", "tags", "first_name", "middle_name", "last_name"})
	//.AddRow(1, time.Now(), time.Now(), "1")
	mock.ExpectQuery(`SELECT string_user FROM %s WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedUserRow)

	mockedLocationRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "status", "tags", "building_numeber", "unit_number", "street_name", "city", "state", "postal_code", "country"})
	mock.ExpectQuery(`SELECT location FROM %s WHERE id = (.+) AND 'deactivated_at' IS NOT NULL`).WithArgs().WillReturnRows(mockedLocationRow)

	instrumentRepo := repository.NewInstrument(sqlxDB)
	userRepo := repository.NewUser(sqlxDB)
	locationRepo := repository.NewLocation(sqlxDB)

	u21Instrument := NewInstrument(instrumentRepo, userRepo, locationRepo)

	u21InstrumentId, err := u21Instrument.Update(instrument)
	assert.NoError(t, err)
	log.Printf("u21InstrumentId: %s", u21InstrumentId)
	assert.Greater(t, len([]rune(u21InstrumentId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new instrument added
	// todo: mock call to client once it's manually tested
}
