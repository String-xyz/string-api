package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/test/stubs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStatus200LoginNoncePayload(t *testing.T) {
	e := echo.New()

	q := make(url.Values)
	q.Set("walletAddress", "walletAddress")

	request := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := NewLogin(nil, stubs.Auth{})
	handler.RegisterRoutes(e.Group("/login"))
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	if assert.NoError(t, handler.NoncePayload(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestStatus200LoginVerifySignature(t *testing.T) {
	e := echo.New()

	body := model.WalletSignaturePayload{}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/sign", strings.NewReader(string(jsonBody)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := NewLogin(nil, stubs.Auth{})
	handler.RegisterRoutes(e.Group("/login"))
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	if assert.NoError(t, handler.VerifySignature(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestStatus400MissingWalletLoginNoncePayload(t *testing.T) {
	e := echo.New()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler := NewLogin(nil, stubs.Auth{})
	handler.RegisterRoutes(e.Group("/login"))
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	if assert.NoError(t, handler.NoncePayload(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
