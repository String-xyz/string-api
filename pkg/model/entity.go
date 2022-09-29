package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Platform struct {
	ID             string         `json:"id" db:"id"`
	CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt  *time.Time     `json:"deactivatedAt" db:"deactivated_at"`
	Type           string         `json:"type" db:"type"`
	ApiKey         string         `json:"apiKey" db:"api_key"`
	Authentication string         `json:"authentication" db:"authentication"`
	Tags           types.JSONText `json:"Tags" db:"tags"`
}

type User struct {
	ID            string         `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt" db:"deactivated_at"`
	Type          string         `json:"type" db:"type"`
	Status        string         `json:"status" db:"status"`
	Tags          types.JSONText `json:"Tags" db:"tags"`
	FirstName     string         `json:"firstName" db:"first_name"`
	MiddleName    string         `json:"middleName" db:"middle_name"`
	LastName      string         `json:"lastName" db:"last_name"`
}

type AuthStrategy struct {
	ID            string     `json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	InvalidatedAt *time.Time `json:"invalidatedAt"`
	AuthType      string     `json:"authType"`
	EntityType    string     `json:"entityType"`
	Token         string     `json:"token"`
}
