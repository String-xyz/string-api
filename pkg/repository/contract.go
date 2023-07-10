package repository

import (
	"context"
	"database/sql"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	"github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"

	"github.com/String-xyz/string-api/pkg/model"
)

type Contract interface {
	database.Transactable
	GetForValidation(ctx context.Context, address string, networkId string, platformId string) (model.Contract, error)
}

type contract[T any] struct {
	repository.Base[T]
}

func NewContract(db database.Queryable) Contract {
	return &contract[model.Contract]{repository.Base[model.Contract]{Store: db, Table: "contract"}}
}

// GetForValidation returns a contract for validation by address, networkId and platformId
func (u contract[T]) GetForValidation(ctx context.Context, address string, networkId string, platformId string) (model.Contract, error) {
	m := model.Contract{}
	err := u.Store.GetContext(ctx, &m, `
	SELECT * FROM contract c 
	INNER JOIN contract_to_platform cp 
	ON c.id = cp.contract_id AND cp.platform_id = $3
	WHERE c.address = $1 AND c.network_id = $2 
	AND deactivated_at IS NULL LIMIT 1
	`,
		address, networkId, platformId)

	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}
