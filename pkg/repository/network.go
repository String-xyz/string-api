package repository

import (
	"context"
	"database/sql"
	"fmt"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	baserepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/model"
)

type Network interface {
	database.Transactable
	Create(ctx context.Context, m model.Network) (model.Network, error)
	GetById(ctx context.Context, id string) (model.Network, error)
	GetByChainId(ctx context.Context, chainId uint64) (model.Network, error)
	Update(ctx context.Context, id string, updates any) error
}

type network[T any] struct {
	baserepo.Base[T]
}

func NewNetwork(db database.Queryable) Network {
	return &network[model.Network]{baserepo.Base[model.Network]{Store: db, Table: "network"}}
}

func (n network[T]) Create(ctx context.Context, insert model.Network) (model.Network, error) {
	m := model.Network{}

	query := `
		INSERT INTO network (name, network_id, chain_id, gas_oracle, rpc_url, explorer_url)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING *`

	row := n.Store.QueryRowxContext(ctx, query, insert.Name, insert.NetworkId, insert.ChainId, insert.GasOracle, insert.RPCUrl, insert.ExplorerUrl)

	err := row.StructScan(&m)
	if err != nil {
		return m, libcommon.StringError(err)
	}

	return m, nil
}

func (n network[T]) GetByChainId(ctx context.Context, chainId uint64) (model.Network, error) {
	m := model.Network{}

	err := n.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1", n.Table), chainId)

	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, nil
}
