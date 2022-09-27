package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type PlaformUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
}

type Platform interface {
	Create(model.Platform) error
	GetID(ID string) (model.Platform, error)
	List(limit int, offset int) ([]model.Platform, error)
	Update(ID string, updates PlaformUpdates) error
}

type platform struct {
	table string
	store *sqlx.DB
}

func NewPlatform(db *sqlx.DB) Platform {
	return &platform{store: db, table: "platform"}
}

func (p platform) Create(m model.Platform) error {
	_, err := p.store.Exec(`
		INSERT INTO platform (type, authentication) 
		VALUES(:type,:authentication)`, m)
	return err
}

func (p platform) GetID(ID string) (model.Platform, error) {
	m := model.Platform{}
	err := p.store.Get(&m, "SELECT FROM platform WHERE id = $1 AND deactivated_at = NULL", ID)
	return m, err
}

func (p platform) List(limit int, offset int) ([]model.Platform, error) {
	list := []model.Platform{}
	if limit == 0 {
		limit = 20
	}
	err := p.store.Select(&list, "SELECT * FROM platform LIMIT $1 OFFSET $2", limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (p platform) Update(ID string, updates PlaformUpdates) error {
	if ID == "" {
		return errors.New("invalid id")
	}
	// Implement updates
	return nil
}
