package api

import (
	"time"

	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/service"
)

func NewRepos(config APIConfig) repository.Repositories {
	// TODO: Make sure all of the repos are initialized here
	return repository.Repositories{
		Auth:           repository.NewAuth(config.Redis, config.DB),
		User:           repository.NewUser(config.DB),
		Contact:        repository.NewContact(config.DB),
		Instrument:     repository.NewInstrument(config.DB),
		Device:         repository.NewDevice(config.DB),
		UserToPlatform: repository.NewUserToPlatform(config.DB),
		Asset:          repository.NewAsset(config.DB),
		Network:        repository.NewNetwork(config.DB),
		Platform:       repository.NewPlatform(config.DB),
		Transaction:    repository.NewTransaction(config.DB),
		TxLeg:          repository.NewTxLeg(config.DB),
		Location:       repository.NewLocation(config.DB),
	}
}

/**
 * Initialize services - Inject Dependencies
 * What dependencies are injectable? DB clients, Store, Repositories, other services, http clients, and every mockable thing
 * Not every service needs access to all of the repos, so we can pass in only the ones it needs. This will make it easier to test
 */
func NewServices(config APIConfig, repos repository.Repositories) service.Services {
	httpClient := service.NewHTTPClient(service.HTTPConfig{Timeout: time.Duration(30) * time.Second})
	client := service.NewFingerprintClient(httpClient)
	fingerprint := service.NewFingerprint(client)
	// we don't need to pass in the entire repos struct, just the ones we need
	verificationRepos := repository.Repositories{Contact: repos.Contact, User: repos.User, Device: repos.Device}
	verification := service.NewVerification(verificationRepos)

	// device service
	deviceRepos := repository.Repositories{Device: repos.Device}
	device := service.NewDevice(deviceRepos, fingerprint)

	auth := service.NewAuth(repos, verification, device)
	apiKey := service.NewAPIKeyStrategy(repos.Auth)
	cost := service.NewCost(config.Redis)
	executor := service.NewExecutor()
	geofencing := service.NewGeofencing(config.Redis)

	// we don't need to pass in the entire repos struct, just the ones we need
	platformRepos := repository.Repositories{Auth: repos.Auth, Platform: repos.Platform}
	platform := service.NewPlatform(platformRepos)

	transaction := service.NewTransaction(repos, config.Redis)
	user := service.NewUser(repos, auth, fingerprint)

	return service.Services{
		Auth:         auth,
		ApiKey:       apiKey,
		Cost:         cost,
		Executor:     executor,
		Geofencing:   geofencing,
		Platform:     platform,
		Transaction:  transaction,
		User:         user,
		Verification: verification,
		Device:       device,
	}
}
