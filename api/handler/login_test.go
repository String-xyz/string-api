package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/String-xyz/string-api/api/validator"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/test/stubs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	testWalletAddress = "0x00196564C0cCb00aC3Eee9F4b564243AF5D69eAc"
	testNonce         = "5df12df88137f80bfce6f05cafa0cd7380f5bd12cb239b29cd65106247ba658ae0f13eace0daef6110ba5be74fe9fbbd19274341f74106aa97d5e04f911323117ad6"
	testSignature     = "3a36eb06f09a1d8d4097101ecf656d11c64bac5bfcd66d0b2325ddabb20f38fdbd1a71e0a04365856c779f6517b4fc2fe8b90d7a93e0b9be9f5b27c3961c150c5f93"
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
	e.Validator = validator.New()
	body := model.WalletSignaturePayload{
		Address:   testNonceWalletAddress,
		Nonce:     testNonce,
		Signature: testSignature,
		Timestamp: time.Now().Unix(),
	}
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
