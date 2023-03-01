package repository

import (
	"database/sql"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Device interface {
	Transactable
	Create(model.Device) (model.Device, error)
	GetById(id string) (model.Device, error)

	// GetByUserIdAndFingerprint gets a device by fingerprint Id and userId, using a compound index
	// the visitor might exisit for two users but the uniqueness comes from (userId, fingerprint)
	GetByUserIdAndFingerprint(userId string, fingerprint string) (model.Device, error)
	GetByUserId(userId string) (model.Device, error)
	ListByUserId(userId string, imit int, offset int) ([]model.Device, error)
	Update(id string, updates any) error
}

type device[T any] struct {
	base[T]
}

func NewDevice(db *sqlx.DB) Device {
	return &device[model.Device]{base[model.Device]{store: db, table: "device"}}
}

func (d device[T]) Create(insert model.Device) (model.Device, error) {
	m := model.Device{}
	rows, err := d.store.NamedQuery(`
		INSERT INTO device (last_used_at,validated_at, type, description, user_id, fingerprint, ip_addresses) 
		VALUES(:last_used_at,:validated_at, :type, :description, :user_id, :fingerprint, :ip_addresses) 
		RETURNING *`, insert)
	if err != nil {
		return m, common.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}

func (d device[T]) GetByUserIdAndFingerprint(userId, fingerprint string) (model.Device, error) {
	m := model.Device{}
	err := d.store.Get(&m, "SELECT * FROM device WHERE user_id = $1 AND fingerprint = $2 LIMIT 1", userId, fingerprint)
	if err != nil && err == sql.ErrNoRows {
		return m, ErrNotFound
	}
	return m, err
}
