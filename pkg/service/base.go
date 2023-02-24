package service

type Services struct {
	Auth   Auth
	ApiKey APIKeyStrategy
	// Chain      Chain    // TODO: Make this service instantiable
	// Checkout   Checkout // TODO: Make this service instantiable
	Cost       Cost
	Executor   Executor
	Geofencing Geofencing
	Platform   Platform
	// Sms          Sms    // TODO: Make this service instantiable
	Transaction  Transaction
	User         User
	Verification Verification
	Device       Device
	Unit21       Unit21
}
