package api

import (
	"time"

	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/service"
)

func NewRepos(config APIConfig) repository.Repositories {
	// TODO: Make sure all of the repos are initialized here
	return repository.Repositories{
		Auth:        repository.NewAuth(config.Redis, config.DB),
		Apikey:      repository.NewApikey(config.DB),
		User:        repository.NewUser(config.DB),
		Contact:     repository.NewContact(config.DB),
		Contract:    repository.NewContract(config.DB),
		Instrument:  repository.NewInstrument(config.DB),
		Device:      repository.NewDevice(config.DB),
		Asset:       repository.NewAsset(config.DB),
		Network:     repository.NewNetwork(config.DB),
		Transaction: repository.NewTransaction(config.DB),
		TxLeg:       repository.NewTxLeg(config.DB),
		Location:    repository.NewLocation(config.DB),
		Platform:    repository.NewPlatform(config.DB),
		Identity:    repository.NewIdentity(config.DB),
	}
}

/**
 * Initialize services - Inject Dependencies
 * What dependencies are injectable? DB clients, Store, Repositories, other services, http clients, and every mockable thing
 * Not every service needs access to all of the repos, so we can pass in only the ones it needs. This will make it easier to test
 */
func NewServices(config APIConfig, repos repository.Repositories) service.Services {
	unit21 := service.NewUnit21(repos)
	httpClient := service.NewHTTPClient(service.HTTPConfig{Timeout: time.Duration(30) * time.Second})
	client := service.NewFingerprintClient(httpClient)
	fingerprint := service.NewFingerprint(client)
	verification := service.NewVerification(repos, unit21)

	// device service
	deviceRepos := repository.Repositories{Device: repos.Device}
	device := service.NewDevice(deviceRepos, fingerprint)

	auth := service.NewAuth(repos, verification, device)
	cost := service.NewCost(config.Redis, repos)
	executor := service.NewExecutor()
	geofencing := service.NewGeofencing(config.Redis)

	transaction := service.NewTransaction(repos, config.Redis, unit21)
	user := service.NewUser(repos, auth, fingerprint, device, unit21, verification)

	card := service.NewCard(repos)

	return service.Services{
		Auth:         auth,
		Cost:         cost,
		Executor:     executor,
		Geofencing:   geofencing,
		Transaction:  transaction,
		User:         user,
		Verification: verification,
		Device:       device,
		Card:         card,
	}
}
