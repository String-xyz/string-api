package stubs

import (
	"context"

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

func (v Verification) SendEmailVerification(ctx context.Context, userId string, email string, platformId string) error {
	return v.Error
}

func (v Verification) VerifyEmail(ctx context.Context, platformId, userId string, deviceId string) error {
	return v.Error
}

func (v Verification) SendDeviceVerification(string, userID string, deviceID string, deviceDescription string) error {
	return v.Error
}

func (v Verification) VerifyDevice(encrypted string) error {
	return v.Error
}

func (v Verification) VerifyEmailWithEncryptedToken(ctx context.Context, encrypted string) error {
	return v.Error
}

func (v Verification) PreValidateEmail(ctx context.Context, platformId, userId, email string) error {
	return v.Error
}

// User Service Stub
type User struct {
	UserOnboardingStatus model.UserOnboardingStatus
	UserCreateResponse   service.UserCreateResponse
	User                 model.User
	Error                error
}

func (u *User) SetOnboardingStatus(m model.UserOnboardingStatus) {
	u.UserOnboardingStatus = m
}

func (u *User) SetResponse(resp service.UserCreateResponse) {
	u.UserCreateResponse = resp
}

func (u *User) SetUser(user model.User) {
	u.User = user
}

func (u User) GetStatus(ctx context.Context, id string) (model.UserOnboardingStatus, error) {
	return u.UserOnboardingStatus, u.Error
}

func (u User) Create(ctx context.Context, request model.WalletSignaturePayloadSigned, platformId string) (service.UserCreateResponse, error) {
	return u.UserCreateResponse, u.Error
}

func (u User) Update(ctx context.Context, userId string, request service.UserUpdates) (model.User, error) {
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

func (a Auth) VerifySignedPayload(ctx context.Context, signature model.WalletSignaturePayloadSigned, platformId string, bypassDevice bool) (service.UserCreateResponse, error) {
	return a.UserCreateResponse, a.Error
}

func (a Auth) GenerateJWT(string, string, ...model.Device) (service.JWT, error) {
	return a.JWT, a.Error
}

func (a Auth) ValidateAPIKeyPublic(key string) (string, error) {
	return "platform-id", a.Error
}

func (a Auth) ValidateAPIKeySecret(key string) (string, error) {
	return "platform-id", a.Error
}

func (a Auth) RefreshToken(ctx context.Context, token string, walletAddress string, platformId string) (service.UserCreateResponse, error) {
	return service.UserCreateResponse{}, a.Error
}

func (a Auth) InvalidateRefreshToken(token string) error {
	return a.Error
}

type Device struct {
	Device model.Device
	Error  error
}

func (d Device) VerifyDevice(ctx context.Context, encrypted string) error {
	return d.Error
}

func (d Device) UpsertDeviceIP(ctx context.Context, deviceId string, ip string) error {
	return d.Error
}

func (d Device) InvalidateUnknownDevice(ctx context.Context, device model.Device) error {
	return d.Error
}

func (d Device) CreateDeviceIfNeeded(ctx context.Context, userId, visitorId, requestId string) (model.Device, error) {
	return d.Device, d.Error
}

func (d Device) CreateUnknownDevice(ctx context.Context, userId string) (model.Device, error) {
	return d.Device, d.Error
}
