package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

func TestBaseUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()
	mock.ExpectExec(`UPDATE contact SET`).WithArgs("type")
	mType := "type"
	m := model.ContactUpdates{Type: &mType}

	NewContact(sqlxDB).Update("Id", m)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error '%s' was not expected, while updating a contact", err)
	}
}
