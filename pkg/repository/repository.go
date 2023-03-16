package repository

type Repositories struct {
	Auth           AuthStrategy
	User           User
	Contact        Contact
	Instrument     Instrument
	Device         Device
	UserToPlatform UserToPlatform
	Asset          Asset
	Network        Network
	Platform       Platform
	Transaction    Transaction
	TxLeg          TxLeg
	Location       Location
}
