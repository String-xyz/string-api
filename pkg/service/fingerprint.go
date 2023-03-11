package service

import (
	"database/sql"
	"errors"

	"github.com/String-xyz/go-lib/common"
	_common "github.com/String-xyz/string-api/pkg/internal/common"
)

type FPClient _common.FingerprintClient
type HTTPConfig _common.HTTPConfig
type HTTPClient _common.HTTPClient
type FPVisitor struct {
	VisitorId  string
	Country    string
	State      string
	IPAddress  sql.NullString
	Timestamp  int64
	Confidence float64
	Type       string
	UserAgent  string
}

func NewHTTPClient(config HTTPConfig) HTTPClient {
	return _common.NewHTTPClient(_common.HTTPConfig(config))
}

func NewFingerprintClient(client HTTPClient) FPClient {
	return _common.NewFingerprint(client)
}

type Fingerprint interface {
	//GetVisitor fetches the visitor data by id, it does not validate if the device is the database
	GetVisitor(id string, request string) (FPVisitor, error)
}

type fingerprint struct {
	client FPClient
}

func NewFingerprint(client FPClient) Fingerprint {
	return &fingerprint{client}
}

func (f fingerprint) GetVisitor(id, requestId string) (FPVisitor, error) {
	visitor, err := f.client.GetVisitorById(id, _common.FPVisitorOpts{Limit: 1, RequestId: requestId})
	if err != nil {
		return FPVisitor{}, common.StringError(err)
	}
	return f.hydrateVisitor(visitor)
}

func (f fingerprint) hydrateVisitor(visitor _common.FPVisitor) (FPVisitor, error) {
	// the check on the lenght here (> 1) is needed since we are always checking the latest visit
	// of the user, if we at some point want to return all the visit, we will need to create a different
	// hydration method.
	if len(visitor.Visits) == 0 || len(visitor.Visits) > 1 {
		return FPVisitor{}, common.StringError(errors.New("visitor history does not match"))
	}

	var state string
	visit := visitor.Visits[0]
	if len(visit.IPLocation.Subdivisions) != 0 {
		state = visit.IPLocation.Subdivisions[0].ISOCode
	}

	return FPVisitor{
		VisitorId:  visitor.Id,
		Country:    visit.IPLocation.Coutry.Code,
		State:      state,
		IPAddress:  sql.NullString{String: visit.IP},
		Timestamp:  visit.Timestamp,
		Confidence: visit.IPLocation.Confidence.Score,
		Type:       visit.BrowserDetails.Device,
		UserAgent:  visit.BrowserDetails.UserAgent,
	}, nil
}
