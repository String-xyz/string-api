package service

import (
	"fmt"
	"net/mail"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/google/uuid"
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

func (a platform) Create(m CreatePlatform) (model.Platform, error) {
	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		return model.Platform{}, errors.Wrap(err, "invalid email")
	}
	tx := a.platRepo.MustBegin()
	a.contactRepo.SetTx(tx)
	defer a.contactRepo.Reset()
	fmt.Println("user id", m.UserID)
	_, err = a.contactRepo.GetUserID(m.UserID)
	if err != nil {
		fmt.Println("err ", err)
		a.platRepo.Rollback()
		return model.Platform{}, err
	}

	plat, err := a.platRepo.Create(m)
	if err != nil {
		a.platRepo.Rollback()
		return model.Platform{}, err
	}

	uuiKey := "str." + uuid.New().String()
	err = a.authRepo.CreateAPIKey(plat.ID, common.ToSha256(uuiKey))
	pt := &plat
	pt.ApiKey = uuiKey
	if err != nil {
		a.platRepo.Rollback()
		return *pt, err
	}

	return plat, a.platRepo.Commit()
}
