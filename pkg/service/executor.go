package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"

	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type ContractCall struct {
	CxAddr     string
	CxABI      []string
	CxFunc     string
	CxParams   []string
	TxValue    int64
	TxGasLimit uint64
}

type CallEstimate struct {
	Value   int64
	Gas     uint64
	Success bool
}

type Executor interface {
	New() Executor
	Initialize(RPC string) (*ethclient.Client, error)
	Estimate(call ContractCall) (CallEstimate, error)
}

type executor struct {
	eth *ethclient.Client
}

func (e executor) New() Executor {
	return &executor{}
}

func NewExecutor() Executor {
	return &executor{}
}

func (e executor) Initialize(RPC string) (*ethclient.Client, error) {
	cl, err := ethclient.Dial(RPC)
	if err != nil {
		return nil, err
	}
	e.eth = cl
	return cl, nil
}

func (e executor) Estimate(call ContractCall) (CallEstimate, error) {
	sk := crypto.ToECDSAUnsafe(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	to := common.HexToAddress(call.CxAddr)
	value := new(big.Int).Mul(big.NewInt(call.TxValue), big.NewInt(params.GWei))
	publicKeyECDSA, ok := sk.Public().(*ecdsa.PublicKey)
	if !ok {
		return CallEstimate{}, errors.New("Estimate: Error casting public key to ECDSA")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasLimit := call.TxGasLimit

	chainId, err := e.eth.ChainID(context.Background())
	if err != nil {
		return CallEstimate{}, err
	}

	nonce, err := e.eth.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return CallEstimate{}, err
	}
	tipCap, _ := e.eth.SuggestGasTipCap(context.Background())
	feeCap, _ := e.eth.SuggestGasPrice(context.Background())

	cost := NewCost(repository.NewCost(nil)) // temporary
	gasGwei, err := cost.QueryOwlracle(chainId.Uint64())
	gasGweiBig := new(big.Int).Mul(big.NewInt(int64(gasGwei)), big.NewInt(params.GWei)) // THIS IS ROUNDING DOWN OUR FLOAT, FIX LATER

	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &to,
			Value:     value,
			Data:      nil,
		},
	)

	msg := ethereum.CallMsg{
		From:      sender,
		To:        &to,
		Gas:       0,
		GasPrice:  gasGweiBig,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Value:     value,
		Data:      nil,
	}

	estimatedGas, err := e.eth.EstimateGas(context.Background(), msg)
	if err != nil {
		return CallEstimate{}, err
	}

	_, failure := e.eth.CallContract(context.Background(), msg, nil)
	success := failure == nil

	// signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainId), sk)
	//return e.eth.SendTransaction(context.Background(), signedTx)
	return CallEstimate{Value: call.TxValue, Gas: estimatedGas, Success: success}, nil
}
