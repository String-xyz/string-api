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
	GetVistor(ID string, requestID string) (FPVisitor, error)
	Update(FPVisitor) error
}

type fingerprint struct {
	client     FPClient
	deviceRepo repository.Device
}

func NewFingerprint(client FPClient, repo repository.Device) Fingerprint {
	return &fingerprint{client, repo}
}

func (f fingerprint) GetVistor(ID, requestID string) (FPVisitor, error) {
	visitor, err := f.client.GetVisitorByID(ID, common.FPVisitorOpts{Limit: 1, RequestID: requestID})
	if err != nil {
		return FPVisitor{}, common.StringError(err)
	}
	return f.hydrateVisitor(visitor)
}

func (f fingerprint) Update(m FPVisitor) error {
	return nil
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
		VisitorID: visitor.ID,
		Country:   visit.IPLocation.Coutry.Code,
		State:     state.ISOCode,
	}, nil
}
