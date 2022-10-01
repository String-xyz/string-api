package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type UserRegister struct {
	FirstNname string `json:"firstName" db:"first_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
	FirstNname    *string         `json:"firstName" db:"first_name"`
	MiddleName    *string         `json:"middleName" db:"middle_name"`
	LastName      *string         `json:"lastName" db:"last_name"`
}
