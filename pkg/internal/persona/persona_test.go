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

func TestGetVerificationById(t *testing.T) {
	testServer := testServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verification := &VerificationResponse{Data: Verification{
			Id:   "test-id",
			Type: "verification",
		}}

		json.NewEncoder(w).Encode(verification)
	}))
	defer func() { testServer.Close() }()

	c := New("test-key")
	c.BaseURL = testServer.URL

	v, err := c.GetVerificationById("test-verification")

	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestCreateInquiry(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/inquiries", r.URL.String())
		assert.Equal(t, http.MethodPost, r.Method)

		var payload InquiryCreateRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "test-account", payload.Data.Attributes.AccountId)
		assert.Equal(t, "test-template", payload.Data.Attributes.TemplateId)

		inquiry := &InquiryResponse{Data: Inquiry{
			Id:   "test-id",
			Type: "inquiry",
		}}

		json.NewEncoder(w).Encode(inquiry)
	})
	defer server.Close()

	client := NewPersonaClient(server.URL, "test-key")
	inquiry, err := client.CreateInquiry(InquiryCreateRequest{
		InquiryCreate{InquiryCreationAttributes{AccountId: "test-account", TemplateId: "test-template"}}})
	assert.NoError(t, err)
	assert.Equal(t, "test-id", inquiry.Data.Id)
	assert.Equal(t, "inquiry", inquiry.Data.Type)
}

func TestGetInquiry(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/inquiries/test-id", r.URL.String())
		assert.Equal(t, http.MethodGet, r.Method)

		inquiry := &InquiryResponse{Data: Inquiry{
			Id:   "test-id",
			Type: "inquiry",
		},
		}

		json.NewEncoder(w).Encode(inquiry)
	})
	defer server.Close()

	client := NewPersonaClient(server.URL, "test-key")
	inquiry, err := client.GetInquiryById("test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", inquiry.Data.Id)
	assert.Equal(t, "inquiry", inquiry.Data.Type)
}
