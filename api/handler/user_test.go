package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	validator "github.com/String-xyz/go-lib/validator"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/test/stubs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStatus200CreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = validator.New()

	body := model.WalletSignaturePayloadSigned{
		Nonce:     testNonce,
		Signature: testSignature,
	}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := NewUser(nil, stubs.User{}, stubs.Verification{})
	handler.RegisterRoutes(e.Group("/users"))
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestStatus200GetUserStatus(t *testing.T) {
	e := echo.New()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetParamNames("id")
	c.SetParamValues("userId")
	c.Set("userId", "userId")

	handler := NewUser(nil, stubs.User{}, stubs.Verification{})
	handler.RegisterRoutes(e.Group("/users"))

	if assert.NoError(t, handler.Status(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestStatus200UserUpdate(t *testing.T) {
	e := echo.New()
	e.Validator = validator.New()

	body := model.UpdateUserName{
		FirstName:  "Testor",
		MiddleName: "Testy",
		LastName:   "Tester",
	}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(jsonBody)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetParamNames("id")
	c.SetParamValues("userId")
	c.Set("userId", "userId")

	handler := NewUser(nil, stubs.User{}, stubs.Verification{})
	handler.RegisterRoutes(e.Group("/users"))

	if assert.NoError(t, handler.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestStatus200VerifyEmail(t *testing.T) {
	e := echo.New()

	q := make(url.Values)
	q.Set("email", "user@email.com")

	request := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetParamNames("id")
	c.SetParamValues("userId")
	c.Set("userId", "userId")

	handler := NewUser(nil, stubs.User{}, stubs.Verification{})
	handler.RegisterRoutes(e.Group("/users"))

	if assert.NoError(t, handler.VerifyEmail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
