package model

import (
	"encoding/json"
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
	Tags          types.JSONText `json:"tags" db:"tags"`
	FirstName     string         `json:"firstName" db:"first_name"`
	MiddleName    string         `json:"middleName" db:"middle_name"`
	LastName      string         `json:"lastName" db:"last_name"`
}

type Contact struct {
	ID                  string     `json:"id" db:"id"`
	UserID              string     `json:"userId" db:"user_id"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
	LastAuthenticatedAt *time.Time `json:"lastAuthenticatedAt" db:"last_authenticated_at"`
	DeactivatedAt       *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type                string     `json:"type" db:"type"`
	Status              string     `json:"status" db:"status"`
	Data                string     `json:"data" db:"data"`
}

type AuthStrategy struct {
	ID            string     `json:"id,omitempty" db:"id"`
	EntityID      string     `json:"entityId" db:"id"` // for redis use only
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"authType" db:"type"`
	EntityType    string     `json:"entityType,omitempty"` // for redis use only
	ContactData   string     `json:"contactData"`          // for redis use only
	ContactID     string     `json:"contactId" db:"contact_id"`
	Data          string     `json:"data" data:"data"`
}

func (a AuthStrategy) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}
