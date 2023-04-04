package repository

import (
	"context"
	"database/sql"
	"fmt"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Contract interface {
	database.Transactable
	Create(ctx context.Context, insert model.Contract) (model.Contract, error)
	GetById(ctx context.Context, id string) (model.Contract, error)
	List(ctx context.Context, limit int, offset int) ([]model.Contract, error)
	Update(ctx context.Context, id string, updates any) error
	GetByAddressAndNetworkAndPlatform(ctx context.Context, address string, networkId string, platformId string) (model.Contract, error)
}

type contract[T any] struct {
	repository.Base[T]
}

func NewContract(db database.Queryable) Contract {
	return &contract[model.Contract]{repository.Base[model.Contract]{Store: db, Table: "contract"}}
}

func (u contract[T]) Create(ctx context.Context, insert model.Contract) (model.Contract, error) {
	m := model.Contract{}
	rows, err := u.Store.NamedQuery(`
		INSERT INTO contract (name, address, functions, network_id, platform_id) 
		VALUES(:name, :address, :functions, :network_id, :platform_id) RETURNING *`, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, libcommon.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}

func (u contract[T]) GetByAddressAndNetworkAndPlatform(ctx context.Context, address string, networkId string, platformId string) (model.Contract, error) {
	m := model.Contract{}
	err := u.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE address = $1 AND network_id = $2 AND platform_id = $3 AND deactivated_at IS NULL LIMIT 1", u.Table), address, networkId, platformId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}
