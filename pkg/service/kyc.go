package service

import (
	"context"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type KYC interface {
	MeetsRequirements(ctx context.Context, userId string, assetType string, cost float64) (met bool, err error)
	GetTransactionLevel(assetType string, cost float64) KYCLevel
	GetUserLevel(ctx context.Context, userId string) (level KYCLevel, err error)
	UpdateUserLevel(ctx context.Context, userId string) (level KYCLevel, err error)
}

type KYCLevel int

const (
	Level0 KYCLevel = iota
	Level1
	Level2
	Level3
)

type kyc struct {
	repos repository.Repositories
}

func NewKYC(repos repository.Repositories) KYC {
	return &kyc{repos}
}

func (k kyc) MeetsRequirements(ctx context.Context, userId string, assetType string, cost float64) (met bool, err error) {
	transactionLevel := k.GetTransactionLevel(assetType, cost)

	userLevel, err := k.GetUserLevel(ctx, userId)
	if err != nil {
		return false, err
	}
	if userLevel >= transactionLevel {
		return true, nil
	} else {
		return false, nil
	}
}

func (k kyc) GetTransactionLevel(assetType string, cost float64) KYCLevel {
	if assetType == "NFT" {
		if cost < 1000.00 {
			return Level1
		} else if cost < 5000.00 {
			return Level2
		} else {
			return Level3
		}
	} else {
		if cost < 5000.00 {
			return Level2
		} else {
			return Level3
		}
	}
}

func (k kyc) GetUserLevel(ctx context.Context, userId string) (level KYCLevel, err error) {
	level, err = k.UpdateUserLevel(ctx, userId)
	if err != nil {
		return level, err
	}

	return level, nil
}

func (k kyc) UpdateUserLevel(ctx context.Context, userId string) (level KYCLevel, err error) {
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
		if identity.Level != int(Level2) {
			identity.Level = int(Level2)
			k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
		}
	} else if points >= 1 && identity.EmailVerified != nil {
		if identity.Level != int(Level1) {
			identity.Level = int(Level1)
			k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
		}
	} else if points <= 1 && identity.Level != int(Level0) {
		identity.Level = int(Level0)
		k.repos.Identity.Update(ctx, userId, model.IdentityUpdates{Level: &identity.Level})
	}

	return KYCLevel(identity.Level), nil
}
