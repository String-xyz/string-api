package model

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type RefreshTokenResponse struct {
	Token string    `json:"token"`
	ExpAt time.Time `json:"expAt"`
}

type SignatureRequest struct {
	Nonce string `json:"nonce" validate:"required,base64"`
}

type JWT struct {
	ExpAt        time.Time            `json:"expAt"`
	IssuedAt     time.Time            `json:"issuedAt"`
	Token        string               `json:"token"`
	RefreshToken RefreshTokenResponse `json:"refreshToken"`
}

type JWTClaims struct {
	UserId     string `json:"userId"`
	PlatformId string `json:"platformId"`
	DeviceId   string `json:"deviceId"`
	jwt.StandardClaims
}

type UserLoginResponse struct {
	JWT  JWT  `json:"authToken"`
	User User `json:"user"`
}

type UserOnboardingStatus struct {
	Status string `json:"status"`
}

type WalletSignaturePayload struct {
	Address   string `json:"address" validate:"required,eth_addr"`
	Timestamp int64  `json:"timestamp" validate:"required"`
}

type FingerprintPayload struct {
	VisitorId string `json:"visitorId"`
	RequestId string `json:"requestId"`
}

type WalletSignaturePayloadSigned struct {
	Nonce       string             `json:"nonce" validate:"required,base64"`
	Signature   string             `json:"signature" validate:"required,base64"`
	Fingerprint FingerprintPayload `json:"fingerprint"`
}
