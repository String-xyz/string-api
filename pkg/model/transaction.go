package model

type TransactionType string

const (
	Unknown    TransactionType = "Unknown"
	Raw        TransactionType = "Raw"
	MintERC721 TransactionType = "MintERC721"
)

type CostEstimate struct {
	Timestamp  int `json:"timestamp" db:"timestamp"`
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

type Transaction struct {
	TxID string `json:"txID"`
}
