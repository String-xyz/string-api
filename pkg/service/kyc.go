package service

import (
	"context"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type KYC interface {
	GetTransactionKYCLevel(assetType string, cost float64) int
	GetUserKYCLevel(ctx context.Context, userId string) (level int, err error)
	UpdateUserKYCLevel(ctx context.Context, userId string) (level int, err error)
}

type kyc struct {
	repos repository.Repositories
}

func NewKYC(repos repository.Repositories) KYC {
	return &kyc{repos}
}

func (k kyc) MeetsRequirements(ctx context.Context, userId string, assetType string, cost float64) (met bool, err error) {
	transactionLevel := k.GetTransactionKYCLevel(assetType, cost)

	userLevel, err := k.GetUserKYCLevel(ctx, userId)
	if err != nil {
		return false, err
	}
	if userLevel >= transactionLevel {
		return true, nil
	} else {
		return false, nil
	}
}

func (k kyc) GetTransactionKYCLevel(assetType string, cost float64) int {
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

func (k kyc) GetUserKYCLevel(ctx context.Context, userId string) (level int, err error) {
	level, err = k.UpdateUserKYCLevel(ctx, userId)
	if err != nil {
		return level, err
	}

	if err != nil {
		return level, err
	}

	return level, nil
}

func (k kyc) UpdateUserKYCLevel(ctx context.Context, userId string) (level int, err error) {
	identity, err := k.repos.Identity.GetByUserId(ctx, userId)
	if err != nil {
		return level, err
	}

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

	if points < 1 && level >= 1 {
		identity.Level = 0
		k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
	} else if points >= 4 && level < 2 {
		identity.Level = 2
		k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
	} else if points >= 1 && level < 1 && identity.EmailVerified != nil {
		identity.Level = 1
		k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
	}

	return identity.Level, nil
}
