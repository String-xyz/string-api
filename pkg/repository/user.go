package repository

import (
	"errors"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type UserUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
	FirstNname    *string         `json:"firstName" db:"first_name"`
	MiddleName    *string         `json:"middleName" db:"middle_name"`
	LastName      *string         `json:"lastName" db:"last_name"`
}

type UserRegister struct {
	FirstNname string `json:"firstName" db:"first_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type User interface {
	Register(UserRegister) error
	Create(model.User) error
	GetID(ID string) (model.User, error)
	List(limit int, offset int) ([]model.User, error)
	Update(ID string, updates UserUpdates) error
}

type user[T any] struct {
	base[T]
}

func NewUser(db *sqlx.DB) User {
	return &user[model.User]{base[model.User]{store: db, table: "string_user"}}
}

// Register registers an user with authentication (email/password)
// this is a rudimentary implementation of onboarding, will later have proper
// onboarding process.
func (u user[T]) Register(UserRegister) error {
	return nil
}

func (u user[T]) Create(m model.User) error {
	_, err := u.store.NamedExec(`
		INSERT INTO string_user (first_name, last_name, type, status) 
		VALUES(:first_name,:last_name, :type, :status)`, m)
	return err
}

func (u user[T]) Update(ID string, updates UserUpdates) error {
	if ID == "" {
		return errors.New("invalid id")
	}
	// Implement updates
	return nil
}
