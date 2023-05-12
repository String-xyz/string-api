package repository

import (
	"context"
	"database/sql"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	baserepo "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Device interface {
	database.Transactable
	Create(ctx context.Context, m model.Device) (model.Device, error)
	GetById(ctx context.Context, id string) (model.Device, error)

	// GetByUserIdAndFingerprint gets a device by fingerprint ID and userId, using a compound index
	// the visitor might exisit for two users but the uniqueness comes from (userId, fingerprint)
	GetByUserIdAndFingerprint(ctx context.Context, userId string, fingerprint string) (model.Device, error)
	GetByUserId(ctx context.Context, id string) (model.Device, error)
	ListByUserId(ctx context.Context, userId string, imit int, offset int) ([]model.Device, error)
	Update(ctx context.Context, id string, updates any) error
}

type device[T any] struct {
	baserepo.Base[T]
}

func NewDevice(db database.Queryable) Device {
	return &device[model.Device]{baserepo.Base[model.Device]{Store: db, Table: "device"}}
}

func (d device[T]) Create(ctx context.Context, insert model.Device) (model.Device, error) {
	m := model.Device{}

	query, args, err := d.Named(`
		INSERT INTO device (last_used_at, validated_at, type, description, user_id, fingerprint, ip_addresses) 
		VALUES(:last_used_at, :validated_at, :type, :description, :user_id, :fingerprint, :ip_addresses) 
		RETURNING *`, insert)

	if err != nil {
		return m, libcommon.StringError(err)
	}

	// Use QueryRowxContext to execute the query with the provided context
	err = d.Store.QueryRowxContext(ctx, query, args...).StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (d device[T]) GetByUserIdAndFingerprint(ctx context.Context, userId, fingerprint string) (model.Device, error) {
	m := model.Device{}
	err := d.Store.GetContext(ctx, &m, "SELECT * FROM device WHERE user_id = $1 AND fingerprint = $2 LIMIT 1", userId, fingerprint)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, err
}
