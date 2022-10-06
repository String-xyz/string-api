package service

import (
	"net/mail"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/pkg/errors"
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
	m := &c
	uuiKey := "str." + uuidWithoutHyphens()
	m.ApiKey = common.ToSha256(uuiKey)
	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		return model.Platform{}, errors.Wrap(err, "invalid email")
	}
	plat, err := a.platRepo.Create(*m)
	if err != nil {
		return model.Platform{}, err
	}

	err = a.authRepo.CreateAPIKey(plat.ID, c.Authentication, common.ToSha256(uuiKey))
	pt := &plat
	pt.ApiKey = uuiKey
	if err != nil {
		return *pt, err
	}

	return plat, nil
}
