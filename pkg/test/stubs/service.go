package stubs

import (
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/service"
)

// Verification Service Stub
type Verification struct {
	Error error
}

func (v *Verification) SetError(e error) {
	v.Error = e
}

func (v Verification) SendEmailVerification(userID string, email string) error {
	return v.Error
}

func (v Verification) VerifyEmail(encrypted string) error {
	return v.Error
}

// User Service Stub
type User struct {
	UserOnboardingStatus model.UserOnboardingStatus
	UserCreateResponse   service.UserCreateResponse
	Error                error
}

func (u *User) SetOnboardinStatus(m model.UserOnboardingStatus) {
	u.UserOnboardingStatus = m
}

func (u *User) SetResponse(resp service.UserCreateResponse) {
	u.UserCreateResponse = resp
}

func (u User) GetStatus(ID string, walletAddress string) (model.UserOnboardingStatus, error) {
	return u.UserOnboardingStatus, u.Error
}

func (u User) Create(request model.WalletSignaturePayload) (service.UserCreateResponse, error) {
	return u.UserCreateResponse, u.Error
}

func (u User) Update(request service.UserUpdates) error {
	return u.Error
}

// Auth Service Stub
type Auth struct {
	WalletSignedPayload model.WalletSignaturePayload
	JWT                 service.JWT
	Error               error
}

func (a *Auth) SetWalletSignedPayload(m model.WalletSignaturePayload) {
	a.WalletSignedPayload = m
}

func (a *Auth) SetJWT(jwt service.JWT) {
	a.JWT = jwt
}

func (a *Auth) SetError(e error) {
	a.Error = e
}

func (a Auth) PayloadToSign(walletAdress string) (model.WalletSignaturePayload, error) {
	return a.WalletSignedPayload, a.Error
}

func (a Auth) VerifySignedPayload(model.WalletSignaturePayload) (service.JWT, error) {
	return a.JWT, a.Error
}

func (a Auth) GenerateJWT(model.User) (service.JWT, error) {
	return a.JWT, a.Error
}

func (a Auth) ValidateAPIKey(key string) bool {
	return true
}
