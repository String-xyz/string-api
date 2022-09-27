package repository

import (
	"database/sql"
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

type User interface {
	Create(model.User) error
	GetID(ID string) (model.User, error)
	List(limit int, offset int) ([]model.User, error)
	Update(ID string, updates UserUpdates) error
}

type user struct {
	table string
	store *sqlx.DB
}

func NewUser(db *sqlx.DB) User {
	return &user{store: db, table: "string-user"}
}

func (u user) Create(m model.User) error {
	_, err := u.store.Exec(`
		INSERT INTO string-user (first_name, last_name, type, status) 
		VALUES(:first_name,:last_name, :type, :status)`, m)
	return err
}

func (u user) GetID(ID string) (model.User, error) {
	m := model.User{}
	err := u.store.Get(&m, "SELECT FROM string-user WHERE id = $1", ID)
	return m, err
}

func (u user) List(limit int, offset int) ([]model.User, error) {
	list := []model.User{}
	if limit == 0 {
		limit = 20
	}
	err := u.store.Select(&list, "SELECT * FROM string-user LIMIT $1 OFFSET $1", limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (u user) Update(ID string, updates UserUpdates) error {
	if ID == "" {
		return errors.New("invalid id")
	}
	// Implement updates
	return nil
}
