package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestKeysAndValues(t *testing.T) {
	mType := "type"
	m := model.UserContactUpdates{Type: &mType}
	names, vals := keysAndValues(m)
	assert.Len(t, names, 1)
	assert.Len(t, vals, 1)
}

func TestBaseUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()
	mock.ExpectExec(`UPDATE contact SET`).WithArgs("type")
	mType := "type"
	m := model.UserContactUpdates{Type: &mType}

	NewUserContact(sqlxDB).Update("ID", m)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error '%s' was not expected, while updating a contact", err)
	}
}
