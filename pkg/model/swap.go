package model

type OffRampRequest struct {
	UserAddress string `json:"userAddress"`
	ChainID     int    `json:"chainID"`
	FromToken   string `json:"fromToken"`
	Amount      string `json:"amount"`
}

type OffRampExecutionRequest struct {
	OffRampRequest
	Quote
	Signature string `json:"signature"`
	CardToken string `json:"cardToken"`
}

type SwapRequest struct {
	UserAddress string `json:"userAddress"`
	ChainID     int    `json:"chainID"`
	FromToken   string `json:"fromToken"`
	ToToken     string `json:"toToken"`
	Amount      string `json:"amount"`
}

type SwapCrossChainRequest struct {
	UserAddress string `json:"userAddress"`
	FromChain   int    `json:"fromChain"`
	ToChain     int    `json:"toChain"`
	FromToken   string `json:"fromToken"`
	ToToken     string `json:"toToken"`
	Amount      string `json:"amount"`
}
