package service

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type CreatePlatform = model.CreatePlatform

type Platform interface {
	Create(CreatePlatform) error
}

type platform struct {
	platRepo    repository.Platform
	contactRepo repository.UserContact
}

func NewPlatform(p repository.Platform, c repository.UserContact) Platform {
	return &platform{p, c}
}

func (a platform) Create(m CreatePlatform) error {
	tx := a.platRepo.MustBegin()
	a.contactRepo.SetTx(tx)
	_, err := a.contactRepo.GetUserID(m.UserID)

	if err != nil {
		a.platRepo.Rollback()
		a.contactRepo.Reset()
		return err
	}

	_, err = a.platRepo.Create(m)
	if err != nil {
		a.platRepo.Rollback()
		a.contactRepo.Reset()
		return err
	}

	return nil
}
