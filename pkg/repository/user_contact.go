package repository

import (
	"database/sql"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type UserContactUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Status        *string    `json:"status" db:"status"`
	Data          *string    `json:"data" db:"data"`
}

type UserContact interface {
	Transactable
	Create(model.Contact) (model.Contact, error)
	GetID(ID string) (model.Contact, error)
	GetUserID(userID string) (model.Contact, error)
	ListUserID(userID string, imit int, offset int) ([]model.Contact, error)
	List(limit int, offset int) ([]model.Contact, error)
	Update(ID string, updates any) error
}

type userContact[T any] struct {
	base[T]
}

func NewUserContact(db *sqlx.DB) UserContact {
	return &userContact[model.Contact]{base: base[model.Contact]{store: db, table: "contact"}}
}

func (u userContact[T]) Create(insert model.Contact) (model.Contact, error) {
	m := model.Contact{}
	rows, err := u.store.NamedQuery(`
		INSERT INTO contact (user_id, data, type, status) 
		VALUES(:user_id, :data, :type, :status) RETURNING *`, insert)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}
	defer rows.Close()
	return m, err
}

func (u userContact[T]) ListUserID(userID string, limit int, offset int) ([]model.Contact, error) {
	list := []model.Contact{}
	if limit == 0 {
		limit = 20
	}
	err := u.store.Select(&list, "SELECT * FROM contact WHERE user_id = $1 LIMIT $2 OFFSET $3", userID, limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}
