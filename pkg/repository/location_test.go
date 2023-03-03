package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	id := uuid.NewString()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "status", "building_number", "unit_number", "street_name", "city", "state", "postal_code", "country"}).
		AddRow(id, time.Now(), time.Now(), "Home", "Verified", "20181", "411", "Lark Avenue", "Somerville", "MA", "01443", "USA")

	mock.ExpectQuery("SELECT * FROM location WHERE id = $1 AND deactivated_at IS NULL").WillReturnRows(rows).WithArgs(id)

	location, err := NewLocation(sqlxDB).GetById(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, location.Id)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error '%s' was not expected, getting location by id", err)
	}
}
