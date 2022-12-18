package service

import (
	"errors"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
)

type FPClient common.FingerprintClient
type HTTPConfig common.HTTPConfig
type HTTPClient common.HTTPClient
type FPVisitor = model.FPVisitor

func NewHTTPClient(config HTTPConfig) HTTPClient {
	return common.NewHTTPClient(common.HTTPConfig(config))
}

func NewFingerprintClient(client HTTPClient) FPClient {
	return common.NewFingerprint(client)
}

type Fingerprint interface {
	//GetVisitor fetches the visitor data by id, it does not validate if the device is the database
	GetVisitor(ID string, request string) (FPVisitor, error)
}

type fingerprint struct {
	client FPClient
}

func NewFingerprint(client FPClient) Fingerprint {
	return &fingerprint{client}
}

func (f fingerprint) GetVisitor(ID, requestID string) (FPVisitor, error) {
	visitor, err := f.client.GetVisitorByID(ID, common.FPVisitorOpts{Limit: 1, RequestID: requestID})
	if err != nil {
		return FPVisitor{}, common.StringError(err)
	}
	return f.hydrateVisitor(visitor)
}

func (f fingerprint) hydrateVisitor(visitor common.FPVisitor) (FPVisitor, error) {
	// the check on the lenght here (> 1) is needed since we are always checking the latest visit
	// of the user, if we at some point want to return all the visit, we will need to create a different
	// hydration method.
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
		IPAddress:  visit.IP,
		Timestamp:  visit.Timestamp,
		Confidence: visit.IPLocation.Confidence.Score,
		OsType:     visit.IPLocation.BrowserDetails.OS,
	}, nil
}
