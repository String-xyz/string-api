package model

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
