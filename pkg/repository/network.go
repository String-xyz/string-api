package repository

import (
	"database/sql"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Network interface {
	Transactable
	Create(model.Network) (model.Network, error)
	GetById(id string) (model.Network, error)
	GetChainID(chainId uint64) (model.Network, error)
	Update(ID string, updates any) error
}

type network[T any] struct {
	base[T]
}

func NewNetwork(db *sqlx.DB) Network {
	return &network[model.Network]{base[model.Network]{store: db, table: "network"}}
}

func (n network[T]) Create(insert model.Network) (model.Network, error) {
	m := model.Network{}
	rows, err := n.store.NamedQuery(`
		INSERT INTO network (name) 
		VALUES(:name) 	RETURNING *`, insert)

	defer rows.Close()

	if err != nil {
		return m, common.StringError(err)
	}

	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, common.StringError(err)
		}
	}

	return m, nil
}

func (n network[T]) GetChainID(chainId uint64) (model.Network, error) {
	m := model.Network{}
	err := n.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1", n.table), chainId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	}
	return m, nil
}
