package service

import (
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Device interface {
	VerifyDevice(encrypted string) error
	CreateDeviceIfNeeded(userID, visitorID, requestID string) (model.Device, error)
	CreateUnknownDevice(userID string) (model.Device, error)
	InvalidateUnknownDevice(device model.Device) error
}

type device struct {
	repos       repository.Repositories
	fingerprint Fingerprint
}

func NewDevice(repos repository.Repositories, f Fingerprint) Device {
	return &device{repos, f}
}

func (d device) createDevice(userID string, visitor FPVisitor, description string) (model.Device, error) {
	addresses := pq.StringArray{}
	if visitor.IPAddress.String != "" {
		addresses = pq.StringArray{visitor.IPAddress.String}
	}

	return d.repos.Device.Create(model.Device{
		UserID:      userID,
		Fingerprint: visitor.VisitorID,
		Type:        visitor.Type,
		IpAddresses: addresses,
		Description: description,
		LastUsedAt:  time.Now(),
	})
}

func (d device) CreateUnknownDevice(userID string) (model.Device, error) {
	visitor := FPVisitor{
		VisitorID: "unknown",
		Type:      "unknown",
		UserAgent: "unknown",
	}
	device, err := d.createDevice(userID, visitor, "an unknown device")
	return device, common.StringError(err)
}

func (d device) CreateDeviceIfNeeded(userID, visitorID, requestID string) (model.Device, error) {
	if visitorID == "" || requestID == "" {
		/* fingerprint is not available, create an unknown device. It should be invalidated on every login */
		device, err := d.getOrCreateUnknownDevice(userID, "unknown")
		if err != nil {
			return device, common.StringError(err)
		}

		if !isDeviceValidated(device) {
			device.ValidatedAt = nil
			return device, nil
		}

		return device, common.StringError(err)
	} else {
		/* device recognized, create or get the device */
		device, err := d.repos.Device.GetByUserIdAndFingerprint(userID, visitorID)
		if err == nil {
			return device, err
		}

		/* create device only if the error is not found */
		if err == repository.ErrNotFound {
			visitor, fpErr := d.fingerprint.GetVisitor(visitorID, requestID)
			if fpErr != nil {
				return model.Device{}, common.StringError(fpErr)
			}
			device, dErr := d.createDevice(userID, visitor, "a new device "+visitor.UserAgent+" ")
			return device, dErr
		}

		return device, common.StringError(err)
	}
}

func (d device) VerifyDevice(encrypted string) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := common.Decrypt[DeviceVerification](encrypted, key)
	if err != nil {
		return common.StringError(err)
	}

	now := time.Now()
	if now.Unix()-received.Timestamp > (60 * 15) {
		return common.StringError(errors.New("link expired"))
	}
	err = d.repos.Device.Update(received.DeviceID, model.DeviceUpdates{ValidatedAt: &now})
	return err
}

func (d device) getOrCreateUnknownDevice(userId, visitorId string) (model.Device, error) {
	var device model.Device

	device, err := d.repos.Device.GetByUserIdAndFingerprint(userId, "unknown")
	if err != nil && err != repository.ErrNotFound {
		return device, common.StringError(err)
	}

	if device.ID != "" {
		return device, nil
	}

	// if device is not found, create a new one
	device, err = d.CreateUnknownDevice(userId)
	return device, common.StringError(err)
}

func isDeviceValidated(device model.Device) bool {
	return device.ValidatedAt != nil && !device.ValidatedAt.IsZero()
}

func (d device) InvalidateUnknownDevice(device model.Device) error {
	if device.Fingerprint != "unknown" {
		return nil // only unknown devices can be invalidated
	}

	device.ValidatedAt = &time.Time{} // Zero time to set it to nil
	return d.repos.Device.Update(device.ID, device)
}
