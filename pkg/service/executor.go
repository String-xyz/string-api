package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"

	str "github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
)

// type ContractCall struct {
// 	CxAddr     string
// 	CxABI      []string
// 	CxFunc     string
// 	CxParams   []string
// 	TxValue    int64
// 	TxGasLimit uint64
// }

// New format
type ContractCall struct {
	CxAddr     string   // Address of contract ie "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	CxFunc     string   // Function declaration ie "mintTo(address) payable"
	CxReturn   string   // Function return type ie "uint256"
	CxParams   []string // Function parameters ie ["0x000000000000000000BEEF", "32"]
	TxValue    string   // Amount of native token to send ie "0.08 ether"
	TxGasLimit uint64   // Gwei gas limit ie 210000
}

type CallEstimate struct {
	Value   int64
	Gas     uint64
	Success bool
}

type Executor interface {
	New() Executor
	Initialize(RPC string) (*w3.Client, error)
	Estimate(call ContractCall) (CallEstimate, error)
}

type executor struct {
	client *w3.Client
	geth   *ethclient.Client
}

func (e executor) New() Executor {
	return &executor{}
}

func NewExecutor() Executor {
	return &executor{}
}

func (e executor) Initialize(RPC string) (*w3.Client, error) {
	cl, err := w3.Dial(RPC)
	if err != nil {
		return nil, err
	}
	// Do it again for our low-level client
	cl2, err := ethclient.Dial(RPC)
	if err != nil {
		return nil, err
	}
	*e.client = *cl
	*e.geth = *cl2
	return cl, nil
}

func (e executor) Estimate(call ContractCall) (CallEstimate, error) {
	sk := crypto.ToECDSAUnsafe(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	to := w3.A(call.CxAddr)
	value := w3.I(call.TxValue)
	publicKeyECDSA, ok := sk.Public().(*ecdsa.PublicKey)
	if !ok {
		return CallEstimate{}, errors.New("Estimate: Error casting public key to ECDSA")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasLimit := call.TxGasLimit

	var chainId64 uint64
	err := e.client.Call(eth.ChainID().Returns(&chainId64))
	if err != nil {
		return CallEstimate{}, err
	}
	// chainId := new(big.Int).SetUint64(chainId64)

	var nonce uint64
	err = e.client.Call(eth.Nonce(sender, nil).Returns(&nonce))
	if err != nil {
		return CallEstimate{}, err
	}

	tipCap, _ := e.geth.SuggestGasTipCap(context.Background())
	feeCap, _ := e.geth.SuggestGasPrice(context.Background())

	cost := NewCost(repository.NewCost(nil)) // temporary
	gasGwei64, err := cost.QueryOwlracle(chainId64)
	if err != nil {
		return CallEstimate{}, err
	}
	gasGwei := new(big.Int).SetUint64(uint64(gasGwei64)) // THIS IS ROUNDING DOWN OUR FLOAT

	var funcEVM = w3.MustNewFunc(call.CxFunc, call.CxReturn)
	data, err := str.ParseParams(funcEVM, call.CxFunc, call.CxParams)
	if err != nil {
		return CallEstimate{}, err
	}

	msg := w3types.Message{
		From:      sender,
		To:        &to,
		Gas:       gasLimit,
		GasPrice:  gasGwei,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Value:     value,
		Input:     data,
	}

	var estimatedGas uint64
	err = e.client.Call(eth.EstimateGas(&msg, nil).Returns(&estimatedGas))
	if err != nil {
		return CallEstimate{}, err
	}

	success := true // maybe EstimateGas will return an error upon evm revert?

	// signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainId), sk)
	//return e.eth.SendTransaction(context.Background(), signedTx)
	return CallEstimate{Value: value.Int64(), Gas: estimatedGas, Success: success}, nil
}
