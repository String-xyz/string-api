package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	m := model.User{
		FirstName: "Mocking",
		LastName:  "Jay",
		Type:      "human",
		Status:    "tested",
	}
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()
	mock.ExpectQuery(`INSERT INTO string_user`).WithArgs(m.FirstName, m.LastName, m.Type, m.Status)

	NewUser(sqlxDB).Create(m)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error '%s' was not expected, while inserting a new user", err)
	}
}

func TestGetUser(t *testing.T) {
	id := uuid.NewString()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "created_at", "updated_at"}).
		AddRow(id, "Mocking", "Jay", time.Now(), time.Now())

	mock.ExpectQuery("SELECT * FROM string_user WHERE id = $1 AND deactivated_at IS NULL").WillReturnRows(rows).WithArgs(id)

	user, err := NewUser(sqlxDB).GetById(id)
	assert.NoError(t, err)
	assert.Equal(t, id, user.ID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error '%s' was not expected, getting user by id", err)
	}
}

func TestListUser(t *testing.T) {
	id1, id2 := uuid.NewString(), uuid.NewString()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "created_at", "updated_at"}).
		AddRow(id1, "Mocking", "Jay", time.Now(), time.Now()).
		AddRow(id2, "Mocking", "Jay", time.Now(), time.Now())

	mock.ExpectQuery("SELECT * FROM string_user LIMIT $1 OFFSET $2").WillReturnRows(rows).WithArgs(10, 0)

	NewUser(sqlxDB).List(10, 0)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error '%s' was not expected, getting the list of users", err)
	}
}
