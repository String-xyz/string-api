package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
)

type Network interface {
	database.Transactable
	Create(model.Network) (model.Network, error)
	GetById(ctx context.Context, id string) (model.Network, error)
	GetByChainId(chainId uint64) (model.Network, error)
	Update(ctx context.Context, id string, updates any) error
}

type network[T any] struct {
	baserepo.Base[T]
}

func NewNetwork(db database.Queryable) Network {
	return &network[model.Network]{baserepo.Base[model.Network]{Store: db, Table: "network"}}
}

func (n network[T]) Create(insert model.Network) (model.Network, error) {
	m := model.Network{}
	rows, err := n.Store.NamedQuery(`
		INSERT INTO network (name, network_id, chain_id, gas_oracle, rpc_url, explorer_url) 
		VALUES(:name, :network_id, :chain_id, :gas_oracle, :rpc_url, :explorer_url) 	RETURNING *`, insert)

	if err != nil {
		return m, common.StringError(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	return m, nil
}

func (n network[T]) GetByChainId(chainId uint64) (model.Network, error) {
	m := model.Network{}
	err := n.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1", n.Table), chainId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, nil
}
