package api

import (
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/service"
)

type services struct {
	Auth service.Auth
	User service.User
	// Device       service.NewDevice
}

func NewRepos(config APIConfig) repository.Repositories {
	return repository.Repositories{
		Auth:         repository.NewAuth(config.Redis, config.DB),
		User:         repository.NewUser(config.DB),
		Contact:      repository.NewContact(config.DB),
		Instrument:   repository.NewInstrument(config.DB),
		Device:       repository.NewDevice(config.DB),
		UserPlatform: repository.NewUserPlatform(config.DB),
		Asset:        repository.NewAsset(config.DB),
		Network:      repository.NewNetwork(config.DB),
		Transaction:  repository.NewTransaction(config.DB),
		TxLeg:        repository.NewTxLeg(config.DB),
		Location:     repository.NewLocation(config.DB),
	}
}

func NewServices(repos repository.Repositories) services {

	Auth := service.NewAuth(repos)
	User := service.NewUser(repos)

	return services{
		Auth: Auth,
		User: User,
	}
}
