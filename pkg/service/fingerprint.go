package service

import (
	"errors"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type FPClient common.FingerprintClient
type HTTPConfig common.HTTPConfig
type HTTPClient common.HTTPClient
type FPVisitor = model.FPVisitor

func NewHTTPClient(config HTTPConfig) HTTPClient {
	return common.NewHTTPClient(common.HTTPConfig(config))
}

func NewFingerPrintClient(client HTTPClient) FPClient {
	return common.NewFingerprint(client)
}

type Fingerprint interface {
	// ValidateDevice uses Fingerprint client to fetch the visitor by a visitorId and requestId
	// this function look up for an associated device as well and attemps to update it
	// Note - the device needs to be associated with an user
	ValidateDevice(ID string, requestID string, userID string) (FPVisitor, error)

	//GetVisitor fetches the visitor data by id, it does not validate if the device is the database
	GetVisitor(ID string, request string) (FPVisitor, error)
}

type fingerprint struct {
	client     FPClient
	deviceRepo repository.Device
}

func NewFingerprint(client FPClient, repo repository.Device) Fingerprint {
	return &fingerprint{client, repo}
}

func (f fingerprint) GetVisitor(ID, requestID string) (FPVisitor, error) {
	visitor, err := f.client.GetVisitorByID(ID, common.FPVisitorOpts{Limit: 1, RequestID: requestID})
	if err != nil {
		return FPVisitor{}, common.StringError(err)
	}

	return f.hydrateVisitor(visitor)
}

func (f fingerprint) ValidateDevice(ID, requestID, userID string) (FPVisitor, error) {
	visitor, err := f.client.GetVisitorByID(ID, common.FPVisitorOpts{Limit: 1, RequestID: requestID})
	if err != nil {
		return FPVisitor{}, common.StringError(err)
	}

	return f.hydrateVisitor(visitor)
}

// Look up the device by fingerprint ID and return a bool indicating if is found
// it disregards the error and threats it has not found.
func (f fingerprint) getDevice(ID string) (model.Device, bool) {
	m, err := f.deviceRepo.GetByFingerprint(ID)
	return m, err == nil
}

func (f fingerprint) updateDevice(m FPVisitor) error {
	return nil
}

func (f fingerprint) validateLocation(device model.Device, visitor FPVisitor) bool {
	return true
}

func (f fingerprint) hydrateVisitor(visitor common.FPVisitor) (FPVisitor, error) {
	if len(visitor.Visits) == 0 || len(visitor.Visits) > 1 {
		return FPVisitor{}, common.StringError(errors.New("visitor history does not match"))
	}

	visit := visitor.Visits[0]
	if len(visit.IPLocation.Subdivisions) == 0 {
		return FPVisitor{}, common.StringError(errors.New("unable to verify user location"))
	}
	state := visit.IPLocation.Subdivisions[0]

	return FPVisitor{
		VisitorID:  visitor.ID,
		Country:    visit.IPLocation.Coutry.Code,
		State:      state.ISOCode,
		Timestamp:  visit.Timestamp,
		Confidence: visit.IPLocation.Confidence.Score,
	}, nil
}
