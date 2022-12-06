package model

type UserOnboardingStatus struct {
	Status string `json:"status"`
}

type WalletSignaturePayload struct {
	Address   string `json:"address" validate:"required,eth_addr"`
	Timestamp int64  `json:"timestamp" validate:"required"`
	Nonce     string `json:"nonce" validate:"required,hexadecimal,len=132"`
	Signature string `json:"signature" validate:"required,hexadecimal,len=132"`
}
