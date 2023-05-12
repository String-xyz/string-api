package service

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/internal/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
	"github.com/pkg/errors"
)

type ContractCall struct {
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
	Initialize(network Chain) error
	Initiate(call ContractCall) (string, *big.Int, error)
	Estimate(call ContractCall) (CallEstimate, error)
	TxWait(txId string) (uint64, error)
	Close() error
	GetByChainId() (uint64, error)
	GetBalance() (float64, error)
}

type executor struct {
	client *w3.Client
	geth   *ethclient.Client
}

func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) Initialize(network Chain) error {
	RPC := network.RPC
	var err error
	e.client, err = w3.Dial(RPC)
	if err != nil {
		return libcommon.StringError(err)
	}
	// Do it again for our low-level client
	e.geth, err = ethclient.Dial(network.RPC)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

func (e *executor) Close() error {
	err := e.client.Close()
	if err != nil {
		return libcommon.StringError(err)
	}
	e.geth.Close()
	return nil
}

func (e executor) Estimate(call ContractCall) (CallEstimate, error) {
	// Generate blockchain message
	msg, err := e.generateTransactionMessage(call)
	if err != nil {
		return CallEstimate{}, libcommon.StringError(err)
	}

	// Estimate gas of message
	var estimatedGas uint64
	err = e.client.Call(eth.EstimateGas(&msg, nil).Returns(&estimatedGas))
	if err != nil {
		// Execution Will Revert!
		return CallEstimate{Value: *msg.Value, Gas: estimatedGas, Success: false}, libcommon.StringError(err)
	}
	return CallEstimate{Value: *msg.Value, Gas: estimatedGas, Success: true}, nil
}

func (e executor) Initiate(call ContractCall) (string, *big.Int, error) {
	tx, err := e.generateTransactionRequest(call)
	if err != nil {
		return "", nil, libcommon.StringError(err)
	}

	// Call tx and retrieve hash
	var hash ethcommon.Hash
	err = e.client.Call(eth.SendTx(&tx).Returns(&hash))
	if err != nil {
		// Execution failed!
		return "", nil, libcommon.StringError(err)
	}
	return hash.String(), tx.Value(), nil
}

func (e executor) TxWait(txId string) (uint64, error) {
	txHash := ethcommon.HexToHash(txId)
	receipt := types.Receipt{}
	for receipt.Status == 0 {
		pendingReceipt, err := e.geth.TransactionReceipt(context.Background(), txHash)
		// TransactionReceipt returns error "not found" while tx is pending
		if err != nil && err.Error() != "not found" {
			return 0, libcommon.StringError(err)
		}
		if pendingReceipt != nil {
			receipt = *pendingReceipt
		}
		// TODO: Sleep for a few ms to keep the cpu cooler
	}
	return receipt.GasUsed, nil
}

func (e executor) GetByChainId() (uint64, error) {
	// Get ChainId from state
	var chainId64 uint64
	err := e.client.Call(eth.ChainID().Returns(&chainId64))
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	return chainId64, nil
}

func (e executor) getAccount() (ethcommon.Address, error) {
	// Get private key
	sk, err := e.getSk()
	if err != nil {
		return ethcommon.Address{}, libcommon.StringError(err)
	}
	// TODO: avoid panicking so that we get an intelligible error message
	publicKeyECDSA, ok := sk.Public().(*ecdsa.PublicKey)
	if !ok {
		return ethcommon.Address{}, libcommon.StringError(errors.New("getAccount: Error casting public key to ECDSA"))
	}
	return ethcommon.HexToAddress(crypto.PubkeyToAddress(*publicKeyECDSA).String()), nil
}

func (e executor) getSk() (ecdsa.PrivateKey, error) {
	// Get private key
	skStr, err := common.DecryptBlobFromKMS(config.Var.EVM_PRIVATE_KEY)
	if err != nil {
		return ecdsa.PrivateKey{}, libcommon.StringError(err)
	}
	sk, err := crypto.ToECDSA(ethcommon.FromHex(skStr))
	if err != nil {
		return ecdsa.PrivateKey{}, libcommon.StringError(err)
	}
	return *sk, nil
}

func (e executor) GetBalance() (float64, error) {
	account, err := e.getAccount()
	if err != nil {
		return 0, libcommon.StringError(err)
	}

	wei := big.Int{}
	err = e.client.Call(eth.Balance(account, nil).Returns(&wei))
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	fwei := new(big.Float)
	fwei.SetString(wei.String())
	balance := new(big.Float).Quo(fwei, big.NewFloat(math.Pow10(18)))
	fbalance, _ := balance.Float64()
	return fbalance, nil // We like thinking in floats
}

func (e executor) generateTransactionMessage(call ContractCall) (w3types.Message, error) {
	sender, err := e.getAccount()
	if err != nil {
		return w3types.Message{}, libcommon.StringError(err)
	}

	to := w3.A(call.CxAddr)
	value := w3.I(call.TxValue)

	// Get ChainId from state
	var chainId64 uint64
	err = e.client.Call(eth.ChainID().Returns(&chainId64))
	if err != nil {
		return w3types.Message{}, libcommon.StringError(err)
	}

	// Get sender nonce
	var nonce uint64
	err = e.client.Call(eth.Nonce(sender, nil).Returns(&nonce))
	if err != nil {
		return w3types.Message{}, libcommon.StringError(err)
	}

	// Get dynamic fee tx gas params
	tipCap, _ := e.geth.SuggestGasTipCap(context.Background())
	feeCap, _ := e.geth.SuggestGasPrice(context.Background())

	// Get handle to function we wish to call
	funcEVM, err := w3.NewFunc(call.CxFunc, call.CxReturn)
	if err != nil {
		return w3types.Message{}, libcommon.StringError(err)
	}

	// Encode function parameters
	data, err := common.ParseEncoding(funcEVM, call.CxFunc, call.CxParams)
	if err != nil {
		return w3types.Message{}, libcommon.StringError(err)
	}

	// Generate blockchain message
	return w3types.Message{
		From:      sender,
		To:        &to,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Value:     value,
		Input:     data,
		Nonce:     nonce,
	}, nil
}

func (e executor) generateTransactionRequest(call ContractCall) (types.Transaction, error) {
	tx := types.Transaction{}

	msg, err := e.generateTransactionMessage(call)
	if err != nil {
		return tx, libcommon.StringError(err)
	}

	// Get chainId from state
	var chainId64 uint64
	err = e.client.Call(eth.ChainID().Returns(&chainId64))
	if err != nil {
		return tx, libcommon.StringError(err)
	}

	// Get dynamic fee tx gas params
	tipCap, _ := e.geth.SuggestGasTipCap(context.Background())
	feeCap, _ := e.geth.SuggestGasPrice(context.Background())

	// Type conversion for chainId
	chainIdBig := new(big.Int).SetUint64(chainId64)

	// Get signer type, this is used to encode the tx
	signer := types.LatestSignerForChainID(chainIdBig)

	// Generate blockchain tx
	dynamicFeeTx := types.DynamicFeeTx{
		ChainID:   chainIdBig,
		Nonce:     msg.Nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       w3.I(call.TxGasLimit).Uint64(),
		To:        msg.To,
		Value:     msg.Value,
		Data:      msg.Input,
	}

	sk, err := e.getSk()
	if err != nil {
		return tx, libcommon.StringError(err)
	}
	// Sign it
	tx = *types.MustSignNewTx(&sk, signer, &dynamicFeeTx)

	return tx, nil
}
