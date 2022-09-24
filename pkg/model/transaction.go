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
	ContractFunction   string   `json:"contractFunction"`   // function declaration, ie "mintTo(address) payable"
	ContractReturn     string   `json:"contractReturn"`     // function return, ie "(uint256)"
	ContractParameters []string `json:"contractParameters"` // All parameters which will be passed into the contractFunction
	TxValue            string   `json:"txValue"`            // cost of transaction ie "0.08 eth"
	GasLimit           string   `json:"gasLimit"`           // maximum gas to be used for transaction ie "21000 gwei"
	Forward            bool     `json:"forward"`            // Forward resulting asset?
}

type Transaction struct {
	TxID string `json:"txID"`
}
