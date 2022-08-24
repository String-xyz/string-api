package model

type Transaction struct {
	ChainID    int    `json:"chainID,omitempty"` // include omitempty if determined necessary
	CreditCard string `json:"creditCard"`
	//Quote Quote `json:"quote"`
}
