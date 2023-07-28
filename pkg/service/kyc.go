package service

import (
	"context"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type KYC interface {
	MeetsRequirements(ctx context.Context, userId string, assetType string, cost float64) (met bool, level int, err error)
	GetTransactionLevel(assetType string, cost float64) int
	GetUserLevel(ctx context.Context, userId string) (level int, err error)
	UpdateUserLevel(ctx context.Context, userId string) (level int, err error)
}

type kyc struct {
	repos repository.Repositories
}

func NewKYC(repos repository.Repositories) KYC {
	return &kyc{repos}
}

func (k kyc) MeetsRequirements(ctx context.Context, userId string, assetType string, cost float64) (met bool, level int, err error) {
	transactionLevel := k.GetTransactionLevel(assetType, cost)

	userLevel, err := k.GetUserLevel(ctx, userId)
	if err != nil {
		return false, transactionLevel, err
	}
	if userLevel >= transactionLevel {
		return true, transactionLevel, nil
	} else {
		return false, transactionLevel, nil
	}
}

func (k kyc) GetTransactionLevel(assetType string, cost float64) int {
	if assetType == "NFT" {
		if cost < 1000.00 {
			return 1
		} else if cost < 5000.00 {
			return 2
		} else {
			return 3
		}
	} else {
		if cost < 5000.00 {
			return 2
		} else {
			return 3
		}
	}
}

func (k kyc) GetUserLevel(ctx context.Context, userId string) (level int, err error) {
	level, err = k.UpdateUserLevel(ctx, userId)
	if err != nil {
		return level, err
	}

	return level, nil
}

func (k kyc) UpdateUserLevel(ctx context.Context, userId string) (level int, err error) {
	identity, err := k.repos.Identity.GetByUserId(ctx, userId)
	if err != nil {
		return level, err
	}

	points := 0
	if identity.EmailVerified != nil {
		points++
	}
	if identity.PhoneVerified != nil {
		points++
	}
	if identity.DocumentVerified != nil {
		points++
	}
	if identity.SelfieVerified != nil {
		points++
	}

	if points >= 4 {
		if identity.Level != 2 {
			identity.Level = 2
			k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
		}
	} else if points >= 1 && identity.EmailVerified != nil {
		if identity.Level != 1 {
			identity.Level = 1
			k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
		}
	} else if points <= 1 && identity.Level != 0 {
		identity.Level = 0
		k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
	}

	return identity.Level, nil
}
