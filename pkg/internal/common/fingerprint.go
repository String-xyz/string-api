package common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type FPVistionBrowserDetails struct {
	BrowserName         string `json:"browserName"`
	BrowserMajorVersion string `json:"browserMajorVersion"`
	BrowserFullVersion  string `json:"browserFullVersion"`
	OS                  string `json:"os"`
	OSVersion           string `json:"osVersion"`
	Device              string `json:"device"`
	UserAgent           string `json:"userAgent"`
}

type FPVisitIpLocation struct {
	AccuracyRadius int     `json:"accuracyRadius"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	PostalCode     string  `json:"postalCode"`
	Timezone       string  `json:"timezone"`
	VisitorFound   bool    `json:"visitorFound"`

	City struct {
		Name string `json:"name"`
	} `json:"city"`

	Coutry struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`

	Continent struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"continent"`

	Subdivisions []struct {
		ISOCode string `json:"isoCode"`
		Name    string `json:"name"`
	} `json:"subdivisions"`
	Confidence struct {
		Score float64 `json:"score"`
	} `json:"confidence"`
}

type FPVisitorVisit struct {
	RequestID      string                  `json:"requestId"`
	Incognito      bool                    `json:"incognito"`
	LinkedId       string                  `json:"linkedId"`
	Time           string                  `json:"time"`
	Timestamp      int64                   `json:"timestamp"`
	URL            string                  `json:"url"`
	IP             string                  `json:"ip"`
	IPLocation     FPVisitIpLocation       `json:"ipLocation"`
	BrowserDetails FPVistionBrowserDetails `json:"browserDetails"`
}

type FPVisitor struct {
	ID     string           `json:"visitorId"`
	Visits []FPVisitorVisit `json:"visits"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPConfig struct {
	Timeout time.Duration
}

type FPVisitorOpts struct {
	Limit     int
	RequestID string
	LinkedID  string
}

func NewHTTPClient(config HTTPConfig) HTTPClient {
	return &http.Client{Timeout: config.Timeout}
}

type FingerprintClient interface {
	// GetVisitorByID get the fingerprint visitor by its id
	// it returns the most up to date information for the  visitor
	// The limit should always be 1 so we can get the latest information
	GetVisitorByID(VisitorID string, opts FPVisitorOpts) (FPVisitor, error)
	Request(method, url string, body io.Reader) (*http.Request, error)
}

type fingerprint struct {
	apiKey  string
	client  HTTPClient
	baseURL string
}

func NewFingerprint(client HTTPClient) FingerprintClient {
	apiKey := os.Getenv("FINGERPRINT_API_KEY")
	baseURL := os.Getenv("FINGERPRINT_API_URL")
	return &fingerprint{client: client, apiKey: apiKey, baseURL: baseURL}
}

func (f fingerprint) GetVisitorByID(visitorID string, opts FPVisitorOpts) (FPVisitor, error) {
	m := FPVisitor{}
	r, err := f.Request(http.MethodGet, f.baseURL+"visitors/"+visitorID, nil)
	if err != nil {
		return m, err
	}
	q := f.optionsToQuery(opts)
	if len(q) > 0 {
		r.URL.RawQuery = q.Encode()
	}
	res, err := f.client.Do(r)
	if err != nil {
		return m, err
	}
	if ok := f.checkStatus(res.StatusCode); !ok {
		return m, errors.New(fmt.Sprintf("fingerprint API call failed with status: %d", res.StatusCode))
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, err
	}
	return parseJSON[FPVisitor](body)
}

func (f fingerprint) checkStatus(status int) bool {
	return status >= 200 && status < 300
}

func (f fingerprint) Request(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Auth-API-Key", f.apiKey)
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, err
}

func (f fingerprint) optionsToQuery(opts FPVisitorOpts) url.Values {
	q := url.Values{}
	if opts.Limit != 0 {
		q.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.RequestID != "" {
		q.Add("request_id", opts.RequestID)
	}
	if opts.LinkedID != "" {
		q.Add("linked_id", opts.LinkedID)
	}
	return q
}
