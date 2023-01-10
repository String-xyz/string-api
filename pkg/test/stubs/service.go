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

func (v Verification) SendDeviceVerification(userID string, deviceID string, deviceDescription string) error {
	return v.Error
}

func (v Verification) VerifyDevice(encrypted string) error {
	return v.Error
}

// User Service Stub
type User struct {
	UserOnboardingStatus model.UserOnboardingStatus
	UserCreateResponse   service.UserCreateResponse
	User                 model.User
	Error                error
}

func (u *User) SetOnboardinStatus(m model.UserOnboardingStatus) {
	u.UserOnboardingStatus = m
}

func (u *User) SetResponse(resp service.UserCreateResponse) {
	u.UserCreateResponse = resp
}

func (u *User) SetUser(user model.User) {
	u.User = user
}

func (u User) GetStatus(ID string) (model.UserOnboardingStatus, error) {
	return u.UserOnboardingStatus, u.Error
}

func (u User) Create(request model.WalletSignaturePayloadSigned) (service.UserCreateResponse, error) {
	return u.UserCreateResponse, u.Error
}

func (u User) Update(userID string, request service.UserUpdates) (model.User, error) {
	return u.User, u.Error
}

// Auth Service Stub
type Auth struct {
	SignablePayload    service.SignablePayload
	UserCreateResponse service.UserCreateResponse
	JWT                service.JWT
	Error              error
}

func (a *Auth) SetWalletSignedPayload(m service.SignablePayload) {
	a.SignablePayload = m
}

func (a *Auth) SetUserCreateResponse(resp service.UserCreateResponse) {
	a.UserCreateResponse = resp
}

func (a *Auth) SetJWT(jwt service.JWT) {
	a.JWT = jwt
}

func (a *Auth) SetError(e error) {
	a.Error = e
}

func (a Auth) PayloadToSign(walletAdress string) (service.SignablePayload, error) {
	return a.SignablePayload, a.Error
}

func (a Auth) VerifySignedPayload(model.WalletSignaturePayloadSigned) (service.UserCreateResponse, error) {
	return a.UserCreateResponse, a.Error
}

func (a Auth) GenerateJWT(model.Device) (service.JWT, error) {
	return a.JWT, a.Error
}

func (a Auth) ValidateAPIKey(key string) bool {
	return true
}

func (a Auth) RefreshToken(token string) (service.JWT, error) {
	return a.JWT, a.Error
}
