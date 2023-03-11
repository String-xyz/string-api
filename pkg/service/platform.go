package service

import (
	"github.com/String-xyz/go-lib/common"
	_common "github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type CreatePlatform = model.CreatePlatform

type Platform interface {
	Create(CreatePlatform) (model.Platform, error)
}

type platform struct {
	repos repository.Repositories
}

func NewPlatform(repos repository.Repositories) Platform {
	return &platform{repos}
}

func (a platform) Create(c CreatePlatform) (model.Platform, error) {
	uuiKey := "str." + uuidWithoutHyphens()
	hashed := _common.ToSha256(uuiKey)
	m := model.Platform{}

	plat, err := a.repos.Platform.Create(m)
	if err != nil {
		return model.Platform{}, common.StringError(err)
	}

	_, err = a.repos.Auth.CreateAPIKey(plat.Id, c.Authentication, hashed, false)
	pt := &plat
	if err != nil {
		return *pt, common.StringError(err)
	}

	return plat, nil
}
