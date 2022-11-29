package model

type UserOnboardingStatus struct {
	Status string `json:"status"`
}

type WalletSignaturePayload struct {
	Address   string `json:"address"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}
