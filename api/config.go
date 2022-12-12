package api

import (
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/service"
)

func NewRepos(config APIConfig) repository.Repositories {
	// TODO: Make sure all of the repos are initialized here
	return repository.Repositories{
		Auth:         repository.NewAuth(config.Redis, config.DB),
		User:         repository.NewUser(config.DB),
		Contact:      repository.NewContact(config.DB),
		Instrument:   repository.NewInstrument(config.DB),
		Device:       repository.NewDevice(config.DB),
		UserPlatform: repository.NewUserPlatform(config.DB),
		Asset:        repository.NewAsset(config.DB),
		Network:      repository.NewNetwork(config.DB),
		Platform:     repository.NewPlatform(config.DB),
		Transaction:  repository.NewTransaction(config.DB),
		TxLeg:        repository.NewTxLeg(config.DB),
		Location:     repository.NewLocation(config.DB),
	}
}

/**
 * Initialize services - Inject Dependencies
 * What dependencies are injectable? DB clients, Store, Repositories, other services, http clients, and every mockable thing
 * Not every service needs access to all of the repos, so we can pass in only the ones it needs. This will make it easier to test
 */
func NewServices(config APIConfig, repos repository.Repositories) service.Services {
	Auth := service.NewAuth(repos)
	ApiKey := service.NewAPIKeyStrategy(repos.Auth)
	Cost := service.NewCost(config.Redis)
	Device := service.NewDeviceService()
	Executor := service.NewExecutor()
	Geofencing := service.NewGeofencing(config.Redis)

	// we don't need to pass in the entire repos struct, just the ones we need
	platformRepos := repository.Repositories{Auth: repos.Auth, Platform: repos.Platform}
	Platform := service.NewPlatform(platformRepos)

	Transaction := service.NewTransaction(repos, config.Redis)
	User := service.NewUser(repos)

	// we don't need to pass in the entire repos struct, just the ones we need
	verificationRepos := repository.Repositories{Contact: repos.Contact, User: repos.User}
	Verification := service.NewVerification(verificationRepos)

	return service.Services{
		Auth:         Auth,
		ApiKey:       ApiKey,
		Cost:         Cost,
		Device:       Device,
		Executor:     Executor,
		Geofencing:   Geofencing,
		Platform:     Platform,
		Transaction:  Transaction,
		User:         User,
		Verification: Verification,
	}
}
