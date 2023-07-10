package model

type TransactionType string

const (
	Unknown    TransactionType = "Unknown"
	Raw        TransactionType = "Raw"
	MintERC721 TransactionType = "MintERC721"
)

type Estimate[T string | float64] struct {
	Timestamp  int64 `json:"timestamp"`
	BaseUSD    T     `json:"baseUSD"`
	GasUSD     T     `json:"gasUSD"`
	TokenUSD   T     `json:"tokenUSD"`
	ServiceUSD T     `json:"serviceUSD"`
	TotalUSD   T     `json:"totalUSD"`
}

type Quote struct {
	TransactionRequest TransactionRequest `json:"request" validate:"required"`
	Estimate           Estimate[string]   `json:"estimate" validate:"required"`
	Signature          string             `json:"signature" validate:"required,base64"`
}

type ExecutionRequest struct {
	Quote       Quote       `json:"quote" validate:"required"`
	PaymentInfo PaymentInfo `json:"paymentInfo" validate:"required"`
}

type PaymentInfo struct {
	CardToken *string `json:"cardToken"`
	CardId    *string `json:"cardId"`
	CVV       *string `json:"cvv"`
	SaveCard  bool    `json:"saveCard" validate:"boolean"`
}

// User will pass this in for a quote and receive Execution Parameters
type TransactionRequest struct {
	UserAddress string   `json:"userAddress" validate:"required,eth_addr"`     // Used to keep track of user ie "0x44A4b9E2A69d86BA382a511f845CbF2E31286770"
	AssetName   string   `json:"assetName" validate:"required,min=3,max=30"`   // Used for receipt
	ChainId     uint64   `json:"chainId" validate:"required,number"`           // Chain ID to execute on e.g. 80000.
	CxAddr      string   `json:"contractAddress" validate:"required,eth_addr"` // Address of contract ie "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	CxFunc      string   `json:"contractFunction" validate:"required"`         // Function declaration ie "mintTo(address)"
	CxReturn    string   `json:"contractReturn"`                               // Function return type ie "uint256"
	CxParams    []string `json:"contractParameters"`                           // Function parameters ie ["0x000000000000000000BEEF", "32"]
	TxValue     string   `json:"txValue"`                                      // Amount of native token to send ie "0.08 ether"
	TxGasLimit  string   `json:"gasLimit" validate:"required,number"`          // Gwei gas limit ie "210000 gwei"
}

type TransactionReceipt struct {
	TxId        string `json:"txId"`
	TxURL       string `json:"txUrl"`
	TxTimestamp string `json:"txTimestamp"`
}
