package service

import (
	"database/sql"
	"errors"

	commonlib "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/pkg/internal/common"
)

type FPClient common.FingerprintClient
type HTTPConfig common.HTTPConfig
type HTTPClient common.HTTPClient
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
	return common.NewHTTPClient(common.HTTPConfig(config))
}

func NewFingerprintClient(client HTTPClient) FPClient {
	return common.NewFingerprint(client)
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
	visitor, err := f.client.GetVisitorById(id, common.FPVisitorOpts{Limit: 1, RequestId: requestId})
	if err != nil {
		return FPVisitor{}, commonlib.StringError(err)
	}
	return f.hydrateVisitor(visitor)
}

func (f fingerprint) hydrateVisitor(visitor common.FPVisitor) (FPVisitor, error) {
	// the check on the lenght here (> 1) is needed since we are always checking the latest visit
	// of the user, if we at some point want to return all the visit, we will need to create a different
	// hydration method.
	if len(visitor.Visits) == 0 || len(visitor.Visits) > 1 {
		return FPVisitor{}, commonlib.StringError(errors.New("visitor history does not match"))
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
