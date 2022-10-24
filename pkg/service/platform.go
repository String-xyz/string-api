package service

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type CreatePlatform = model.CreatePlatform

type Platform interface {
	Create(CreatePlatform) (model.Platform, error)
}

type platform struct {
	platRepo    repository.Platform
	contactRepo repository.UserContact
	authRepo    repository.AuthStrategy
}

func NewPlatform(p repository.Platform, c repository.UserContact, a repository.AuthStrategy) Platform {
	return &platform{p, c, a}
}

func (a platform) Create(c CreatePlatform) (model.Platform, error) {
	uuiKey := "str." + uuidWithoutHyphens()
	hashed := common.ToSha256(uuiKey)
	m := model.Platform{
		Type:           c.Type,
		Authentication: c.Authentication,
		ApiKey:         hashed,
		Status:         "initialized",
	}

	plat, err := a.platRepo.Create(m)
	if err != nil {
		return model.Platform{}, err
	}

	err = a.authRepo.CreateAPIKey(plat.ID, c.Authentication, hashed)
	pt := &plat
	pt.ApiKey = uuiKey
	if err != nil {
		return *pt, err
	}

	return plat, nil
}
