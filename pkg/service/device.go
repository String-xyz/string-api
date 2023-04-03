package service

import (
	"context"
	"os"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/internal/common"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"

	"github.com/lib/pq"
)

type Device interface {
	VerifyDevice(ctx context.Context, encrypted string) error
	UpsertDeviceIP(ctx context.Context, deviceId string, Ip string) (err error)
	InvalidateUnknownDevice(ctx context.Context, device model.Device) error
	CreateDeviceIfNeeded(userId, visitorId, requestId string) (model.Device, error)
	CreateUnknownDevice(userId string) (model.Device, error)
}

type device struct {
	repos       repository.Repositories
	fingerprint Fingerprint
}

func NewDevice(repos repository.Repositories, f Fingerprint) Device {
	return &device{repos, f}
}

func (d device) VerifyDevice(ctx context.Context, encrypted string) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := libcommon.Decrypt[DeviceVerification](encrypted, key)
	if err != nil {
		return libcommon.StringError(err)
	}

	now := time.Now()
	if now.Unix()-received.Timestamp > (60 * 15) {
		return libcommon.StringError(serror.EXPIRED)
	}
	err = d.repos.Device.Update(ctx, received.DeviceId, model.DeviceUpdates{ValidatedAt: &now})
	return err
}

func (d device) UpsertDeviceIP(ctx context.Context, deviceId string, ip string) (err error) {
	device, err := d.repos.Device.GetById(ctx, deviceId)
	if err != nil {
		return
	}
	contains := common.SliceContains(device.IpAddresses, ip)
	if !contains {
		ipAddresses := append(device.IpAddresses, ip)
		updates := &model.DeviceUpdates{IpAddresses: &ipAddresses}
		err = d.repos.Device.Update(ctx, deviceId, updates)
		if err != nil {
			return
		}
	}
	return
}

func (d device) CreateDeviceIfNeeded(userId, visitorId, requestId string) (model.Device, error) {
	if visitorId == "" || requestId == "" {
		/* fingerprint is not available, create an unknown device. It should be invalidated on every login */
		device, err := d.getOrCreateUnknownDevice(userId, "unknown")
		if err != nil {
			return device, libcommon.StringError(err)
		}

		if !isDeviceValidated(device) {
			device.ValidatedAt = nil
			return device, nil
		}

		return device, libcommon.StringError(err)
	} else {
		/* device recognized, create or get the device */
		device, err := d.repos.Device.GetByUserIdAndFingerprint(userId, visitorId)
		if err == nil {
			return device, err
		}

		/* create device only if the error is not found */
		if serror.Is(err, serror.NOT_FOUND) {
			visitor, fpErr := d.fingerprint.GetVisitor(visitorId, requestId)
			if fpErr != nil {
				return model.Device{}, libcommon.StringError(fpErr)
			}
			device, dErr := d.createDevice(userId, visitor, "a new device "+visitor.UserAgent+" ")
			return device, dErr
		}

		return device, libcommon.StringError(err)
	}
}

func (d device) CreateUnknownDevice(userId string) (model.Device, error) {
	visitor := FPVisitor{
		VisitorId: "unknown",
		Type:      "unknown",
		UserAgent: "unknown",
	}
	device, err := d.createDevice(userId, visitor, "an unknown device")
	return device, libcommon.StringError(err)
}

func (d device) InvalidateUnknownDevice(ctx context.Context, device model.Device) error {
	if device.Fingerprint != "unknown" {
		return nil // only unknown devices can be invalidated
	}

	device.ValidatedAt = &time.Time{} // Zero time to set it to nil
	return d.repos.Device.Update(ctx, device.Id, device)
}

func (d device) createDevice(userId string, visitor FPVisitor, description string) (model.Device, error) {
	addresses := pq.StringArray{}
	if visitor.IPAddress.String != "" {
		addresses = pq.StringArray{visitor.IPAddress.String}
	}

	return d.repos.Device.Create(model.Device{
		UserId:      userId,
		Fingerprint: visitor.VisitorId,
		Type:        visitor.Type,
		IpAddresses: addresses,
		Description: description,
		LastUsedAt:  time.Now(),
	})
}

func (d device) getOrCreateUnknownDevice(userId, visitorId string) (model.Device, error) {
	var device model.Device

	device, err := d.repos.Device.GetByUserIdAndFingerprint(userId, "unknown")
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return device, libcommon.StringError(err)
	}

	if device.Id != "" {
		return device, nil
	}

	// if device is not found, create a new one
	device, err = d.CreateUnknownDevice(userId)
	return device, libcommon.StringError(err)
}

func isDeviceValidated(device model.Device) bool {
	return device.ValidatedAt != nil && !device.ValidatedAt.IsZero()
}
