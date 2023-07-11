package persona

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestNew(t *testing.T) {
	c := New("test-key")
	if c == nil {
		t.Error("Expected persona client to be created")
	}
}

func TestDoRequest(t *testing.T) {
	testServer := testServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`OK`))
	}))
	defer func() { testServer.Close() }()

	c := New("test-key")
	c.BaseURL = testServer.URL

	err := c.doRequest("GET", "/", nil, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetVerifications(t *testing.T) {
	testServer := testServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":"test-verification"}]`))
	}))
	defer func() { testServer.Close() }()

	c := New("test-key")
	c.BaseURL = testServer.URL

	v, err := c.GetVerifications()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(v) != 1 {
		t.Errorf("Expected one verification, got %v", len(v))
	}
}

func TestGetTemplates(t *testing.T) {
	testServer := testServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":"test-template"}]`))
	}))
	defer func() { testServer.Close() }()

	c := New("test-key")
	c.BaseURL = testServer.URL

	te, err := c.GetTemplates()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(te) != 1 {
		t.Errorf("Expected one template, got %v", len(te))
	}
}

func TestCreateInquiry(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/inquiries", r.URL.String())
		assert.Equal(t, http.MethodPost, r.Method)

		var payload InquiryPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "test-account", payload.AccountID)
		assert.Equal(t, "test-template", payload.Template)

		inquiry := &Inquiry{
			Id:     "test-id",
			Status: "completed",
		}

		json.NewEncoder(w).Encode(inquiry)
	})
	defer server.Close()

	client := NewPersonaClient(server.URL, "test-key")
	inquiry, err := client.CreateInquiry(InquiryPayload{
		AccountID: "test-account",
		Template:  "test-template",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-id", inquiry.Id)
	assert.Equal(t, "completed", inquiry.Status)
}

func TestGetInquiry(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/inquiries/test-id", r.URL.String())
		assert.Equal(t, http.MethodGet, r.Method)

		inquiry := &Inquiry{
			Id:     "test-id",
			Status: "completed",
		}

		json.NewEncoder(w).Encode(inquiry)
	})
	defer server.Close()

	client := NewPersonaClient(server.URL, "test-key")
	inquiry, err := client.GetInquiry("test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", inquiry.Id)
	assert.Equal(t, "completed", inquiry.Status)
}
