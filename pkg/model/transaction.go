package model

type TransactionType string

const (
	Unknown    TransactionType = "Unknown"
	Raw        TransactionType = "Raw"
	MintERC721 TransactionType = "MintERC721"
)

type Quote struct {
	Timestamp  int64   `json:"timestamp"`
	BaseUSD    float64 `json:"baseUSD"`
	GasUSD     float64 `json:"gasUSD"`
	TokenUSD   float64 `json:"tokenUSD"`
	ServiceUSD float64 `json:"serviceUSD"`
	TotalUSD   float64 `json:"totalUSD"`
}

type ExecutionRequest struct {
	TransactionRequest
	Quote
	Signature string `json:"signature"`
	CardToken string `json:"cardToken"`
}

type PrecisionSafeQuote struct {
	Timestamp  int64  `json:"timestamp"`
	BaseUSD    string `json:"baseUSD"`
	GasUSD     string `json:"gasUSD"`
	TokenUSD   string `json:"tokenUSD"`
	ServiceUSD string `json:"serviceUSD"`
	TotalUSD   string `json:"totalUSD"`
}

type PrecisionSafeExecutionRequest struct {
	TransactionRequest
	PrecisionSafeQuote
	Signature    string `json:"signature"`
	CardToken    string `json:"cardToken"`
	CardSourceId string `json:"cardSourceId"`
}

// User will pass this in for a quote and receive Execution Parameters
type TransactionRequest struct {
	UserAddress string   `json:"userAddress"`        // Used to keep track of user ie "0x44A4b9E2A69d86BA382a511f845CbF2E31286770"
	ChainId     uint64   `json:"chainId"`            // Chain ID to execute on e.g. 80000
	CxAddr      string   `json:"contractAddress"`    // Address of contract ie "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	CxFunc      string   `json:"contractFunction"`   // Function declaration ie "mintTo(address) payable"
	CxReturn    string   `json:"contractReturn"`     // Function return type ie "uint256"
	CxParams    []string `json:"contractParameters"` // Function parameters ie ["0x000000000000000000BEEF", "32"]
	TxValue     string   `json:"txValue"`            // Amount of native token to send ie "0.08 ether"
	TxGasLimit  string   `json:"gasLimit"`           // Gwei gas limit ie "210000 gwei"
}

type TransactionReceipt struct {
	TxId  string `json:"txId"`
	TxURL string `json:"txUrl"`
}
