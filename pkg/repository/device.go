package repository

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Device interface {
	Transactable
	Create(model.Device) (model.Device, error)
	GetID(id string) (model.Device, error)
	Update(ID string, updates any) error
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
		INSERT INTO device (name) 
		VALUES(:name) 	RETURNING *`, insert)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		err = rows.StructScan(&m)
	}

	defer rows.Close()
	return m, err
}
