package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"

	stringCommon "github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
)

type ContractCall struct {
	RPC        string
	CxAddr     string   // Address of contract ie "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	CxFunc     string   // Function declaration ie "mintTo(address) payable"
	CxReturn   string   // Function return type ie "uint256"
	CxParams   []string // Function parameters ie ["0x000000000000000000BEEF", "32"]
	TxValue    string   // Amount of native token to send ie "0.08 ether"
	TxGasLimit string   // Gwei gas limit ie "210000 gwei"
}

type CallEstimate struct {
	Value   big.Int
	Gas     uint64
	Success bool
}

type Executor interface {
	New() Executor
	Initialize(RPC string) error
	Initiate(call ContractCall) (string, error)
	Estimate(call ContractCall) (CallEstimate, error)
	Close() error
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

func (e *executor) Initialize(RPC string) error {
	var err error
	e.client, err = w3.Dial(RPC)
	if err != nil {
		return err
	}
	// Do it again for our low-level client
	e.geth, err = ethclient.Dial(RPC)
	if err != nil {
		return err
	}
	return nil
}

func (e *executor) Close() error {
	err := e.client.Close()
	if err != nil {
		return err
	}
	e.geth.Close()
	return nil
}

func (e executor) Estimate(call ContractCall) (CallEstimate, error) {
	// Get private key
	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	if err != nil {
		return CallEstimate{}, err
	}
	// TODO: avoid panicking so that we get an intelligible error message
	to := w3.A(call.CxAddr)
	value := w3.I(call.TxValue)
	// Get public key
	publicKeyECDSA, ok := sk.Public().(*ecdsa.PublicKey)
	if !ok {
		return CallEstimate{}, errors.New("Estimate: Error casting public key to ECDSA")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get ChainID from state
	var chainId64 uint64
	err = e.client.Call(eth.ChainID().Returns(&chainId64))
	if err != nil {
		return CallEstimate{}, err
	}

	// Get sender nonce
	var nonce uint64
	err = e.client.Call(eth.Nonce(sender, nil).Returns(&nonce))
	if err != nil {
		return CallEstimate{}, err
	}

	// Get dynamic fee tx gas params
	tipCap, _ := e.geth.SuggestGasTipCap(context.Background())
	feeCap, _ := e.geth.SuggestGasPrice(context.Background())

	// Get handle to function we wish to call
	funcEVM, err := w3.NewFunc(call.CxFunc, call.CxReturn)
	if err != nil {
		return CallEstimate{}, err
	}

	// Encode function parameters
	data, err := stringCommon.ParseEncoding(funcEVM, call.CxFunc, call.CxParams)
	if err != nil {
		return CallEstimate{}, err
	}

	// Generate blockchain message
	msg := w3types.Message{
		From:      sender,
		To:        &to,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Value:     value,
		Input:     data,
	}

	// Estimate gas of message
	var estimatedGas uint64
	err = e.client.Call(eth.EstimateGas(&msg, nil).Returns(&estimatedGas))
	if err != nil {
		// Execution Will Revert!
		return CallEstimate{Value: *value, Gas: estimatedGas, Success: false}, nil
	}
	return CallEstimate{Value: *value, Gas: estimatedGas, Success: true}, nil
}

func (e executor) Initiate(call ContractCall) (string, error) {
	// Get private key
	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	if err != nil {
		return "", err
	}
	// TODO: avoid panicking so that we get an intelligible error message
	to := w3.A(call.CxAddr)
	value := w3.I(call.TxValue)
	// Get public key
	publicKeyECDSA, ok := sk.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("Estimate: Error casting public key to ECDSA")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Use provided gas limit
	gasLimit := w3.I(call.TxGasLimit)

	// Get chainID from state
	var chainId64 uint64
	err = e.client.Call(eth.ChainID().Returns(&chainId64))
	if err != nil {
		return "", err
	}

	// Get sender nonce
	var nonce uint64
	err = e.client.Call(eth.Nonce(sender, nil).Returns(&nonce))
	if err != nil {
		return "", err
	}

	// Get dynamic fee tx gas params
	tipCap, _ := e.geth.SuggestGasTipCap(context.Background())
	feeCap, _ := e.geth.SuggestGasPrice(context.Background())

	// Get handle to function we wish to call
	funcEVM, err := w3.NewFunc(call.CxFunc, call.CxReturn)
	if err != nil {
		return "", err
	}

	// Encode function parameters
	data, err := stringCommon.ParseEncoding(funcEVM, call.CxFunc, call.CxParams)
	if err != nil {
		return "", err
	}

	// Type conversion for chainID
	chainIdBig := new(big.Int).SetUint64(chainId64)

	// Get signer type, this is used to encode the tx
	signer := types.LatestSignerForChainID(chainIdBig)

	// Generate blockchain tx
	dynamicFeeTx := types.DynamicFeeTx{
		ChainID:   chainIdBig,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit.Uint64(),
		To:        &to,
		Value:     value,
		Data:      data,
	}
	// Sign it
	tx := types.MustSignNewTx(sk, signer, &dynamicFeeTx)

	// Call tx and retrieve hash
	var hash common.Hash
	err = e.client.Call(eth.SendTx(tx).Returns(&hash))
	if err != nil {
		// Execution failed!
		return "", err
	}
	return hash.String(), nil
}

func (e executor) TxWait(txID string) (uint64, error) {
	txHash := common.HexToHash(txID)
	receipt := types.Receipt{}
	for receipt.Status == 0 {
		pendingReceipt, err := e.geth.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return 0, err
		}
		if pendingReceipt != nil {
			receipt = *pendingReceipt
		}
		// TODO: Sleep for a few ms to keep the cpu cooler
	}
	return receipt.GasUsed, nil
}
