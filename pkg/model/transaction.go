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
	UserAddress        string   `json:"userAddress"`        // users wallet
	ContractAddress    string   `json:"contractAddress"`    // 0x0000 or ENS name for contract
	ContractABI        []string `json:"contractABI"`        // relevant declarations of contract ABI
	ContractFunction   string   `json:"contractFunction"`   // function name, i.e. 'transfer' or 'mint'
	ContractParameters []string `json:"contractParameters"` // All parameters which will be passed into the contractFunction
	TxValue            string   `json:"txValue"`            // gwei cost of transaction
	GasLimit           string   `json:"gasLimit"`           // maximum gas to be used for transaction
	Forward            bool     `json:"forward"`            // Forward resulting asset?
}

type Transaction struct {
	TxID string `json:"txID"`
}
