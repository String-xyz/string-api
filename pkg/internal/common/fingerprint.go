package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type FPVistionBrowerDetails struct {
	BrowerName           string `json:"browserName"`
	BroweserMajorVersion string `json:"browserMajorVersion"`
	BrowserFullVersion   string `json:"browserFullVersion"`
	OS                   string `json:"os"`
	OSVersion            string `json:"osVersion"`
	Device               string `json:"device"`
	UserAgent            string `json:"userAgent"`
}

type FPVisitIpLocation struct {
	AccuracyRadius int                    `json:"accuracyRadius"`
	Latitude       float64                `json:"latitude"`
	Longitude      float64                `json:"longitude"`
	PostalCode     string                 `json:"postalCode"`
	Timezone       string                 `json:"timezone"`
	VisitorFound   bool                   `json:"visitorFound"`
	BrowserDetails FPVistionBrowerDetails `json:"browserDetails"`
	City           struct {
		Name string `json:"name"`
	} `json:"city"`

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
	RequestID  string            `json:"requestId"`
	Incognito  bool              `json:"incognito"`
	LinkedId   string            `json:"linkedId"`
	Time       time.Time         `json:"time"`
	Timestamp  time.Time         `json:"timestamp"`
	URL        string            `json:"url"`
	IP         string            `json:"ip"`
	IPLocation FPVisitIpLocation `json:"ipLocation"`
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

type Fingerprint struct {
	apiKey  string
	client  HTTPClient
	baseURL string
}

func NewHTTPClient(config HTTPConfig) HTTPClient {
	return &http.Client{Timeout: config.Timeout}
}

func NewFingerprint(client HTTPClient) *Fingerprint {
	apiKey := os.Getenv("FINGERPRINT_API_KEY")
	return &Fingerprint{client: client, apiKey: apiKey, baseURL: "https://api.fpjs.io/"}
}

func (f Fingerprint) GetVisitorByID(visitonID string) (FPVisitor, error) {
	m := FPVisitor{}
	r, err := http.NewRequest(http.MethodGet, f.pathWithKey("visitors/"+visitonID), nil)
	if err != nil {
		return m, err
	}
	r.Header.Add("Accept", "application/json")
	res, err := f.client.Do(r)
	if err != nil {
		return m, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, err
	}

	return parseJSON[FPVisitor](body)
}

func (f Fingerprint) pathWithKey(path string) string {
	return fmt.Sprintf("%s%s?api_key=%s", f.baseURL, path, f.apiKey)
}

func parseJSON[T any](b []byte) (T, error) {
	var r T
	if err := json.Unmarshal(b, &r); err != nil {
		return r, err
	}
	return r, nil
}

func toJSON(T any) ([]byte, error) {
	return json.Marshal(T)
}
