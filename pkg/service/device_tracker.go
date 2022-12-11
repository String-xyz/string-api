package service

type Device interface {
	IsDeviceAllowed(visitorId string) (bool, error)
	RegisterNewUserDevice(string, string) error
}

type device struct {
}

func NewDeviceService() Device {
	return &device{}
}

func (d device) IsDeviceAllowed(visitorId string) (bool, error) {
	return true, nil
}

func (d device) RegisterNewUserDevice(userId string, deviceId string) error {
	return nil
}
