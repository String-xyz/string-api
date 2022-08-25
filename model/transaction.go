package model

type TransactionType int8

const (
	Unknown TransactionType = iota
	Raw
	MintERC721
)

type CostEstimate struct {
	Timestamp  int
	BaseUSD    int
	GasUSD     int
	TokenUSD   int
	ServiceUSD int
	TotalUSD   int
}

type SignedQuote struct {
	Estimate  CostEstimate
	Signature string
}

type TransactionRequest struct {
	TransactionType
	SignedQuote
	CardToken string `json:"cardToken"`
	ChainID   int    `json:"chainID"`
}

type TransactionData struct {
	TransactionRequest
	UserAddress        string   // users wallet
	ContractAddress    string   // 0x0000 or ENS name for contract
	ContractABI        []string // relevant declarations of contract ABI
	ContractFunction   string   // function name, i.e. 'transfer' or 'mint'
	ContractParameters []string // All parameters which will be passed into the contractFunction
	TxValue            string
	GasLimit           string
	Forward            bool // Forward resulting asset?
}

type TransactionResponse struct {
	TxID string `json:"txID"`
}
