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

func TestCreateEntity(t *testing.T) {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	id := uuid.NewString()
	db, _, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("error %s was not expected when opening stub db", err)
	}
	defer db.Close()

	user := model.User{
		ID:            id,
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

	// Dependent on Device and Instrument Repos being created
	userRepo := repository.NewUser(sqlxDB)
	// deviceRepo := repository.NewDevice(sqlxDB)
	contactRepo := repository.NewUserContact(sqlxDB)
	// instrumentRepo := repository.NewInstrument(sqlxDB)

	// u21Entity := newEntity(userRepo, deviceRepo, contactRepo, instrumentRepo)
	u21Entity := newEntity(userRepo, contactRepo)

	u21EntityId, err := u21Entity.Create(user)
	assert.NoError(t, err)
	log.Printf("u21EntityId: %s", u21EntityId)
	assert.Greater(t, len([]rune(u21EntityId)), 0)

	//validate response from Unit21
	//check Unit21 dashboard for new entity added
	// todo: mock call to client once it's manually tested
}
