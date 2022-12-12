package model

type UserOnboardingStatus struct {
	Status string `json:"status"`
}

type WalletSignaturePayload struct {
	Address   string `json:"address" validate:"required,eth_addr"`
	Timestamp int64  `json:"timestamp" validate:"required"`
}

type WalletSignaturePayloadSigned struct {
	Nonce     string `json:"nonce" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}
