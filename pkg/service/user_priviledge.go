package service

import (
	"context"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/pkg/repository"
)

func UserHasPriviledge(ctx context.Context, cost float64, assetType string, userAddress string, repos repository.Repositories) (bool, int, error) {
	instrument, err := repos.Instrument.GetWalletByAddr(ctx, userAddress)
	if err != nil {
		return false, -1, common.StringError(err)
	}
	_ /*user*/, err = repos.User.GetById(ctx, instrument.UserId)
	if err != nil {
		return false, -1, common.StringError(err)
	}
	priviledge := 3 // TODO: Get this from user once it's in the table

	requirement := 0
	if assetType == "NFT" {
		if cost < 1000.00 {
			requirement = 1
		} else if cost < 5000.00 {
			requirement = 2
		} else {
			requirement = 3
		}
	} else {
		if cost < 5000.00 {
			requirement = 2
		} else {
			requirement = 3
		}
	}
	if priviledge >= requirement {
		return true, requirement, nil
	}
	return false, requirement, nil
}
