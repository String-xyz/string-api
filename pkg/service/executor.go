package service

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"strings"

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
	Initiate(calls []ContractCall) ([]string, *big.Int, error)
	Estimate(calls []ContractCall) (CallEstimate, error)
	TxWait(txIds []string) (uint64, error)
	Close() error
	GetByChainId() (uint64, error)
	GetBalance() (float64, error)
	GetTokenIds(txIds []string) ([]string, error)
	GetTokenQuantities(txIds []string) ([]string, error)
	GetEventData(txIds []string, eventSignature string) ([]types.Log, error)
	ForwardNonFungibleTokens(txIds []string, recipient string) ([]string, []string, error)
	ForwardTokens(txIds []string, recipient string) ([]string, []string, []string, error)
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

func (e executor) Estimate(calls []ContractCall) (CallEstimate, error) {
	// Generate blockchain messages
	w3calls := []w3types.Caller{}
	estimatedGasses := make([]uint64, len(calls))
	totalValue := big.NewInt(0)
	includesApprove := false
	for i, call := range calls {
		// TODO: Further optimize this to take in the calls array
		msg, err := e.generateTransactionMessage(call, uint64(i))
		if err != nil {
			return CallEstimate{}, libcommon.StringError(err)
		}
		totalValue = totalValue.Add(totalValue, msg.Value)

		w3calls = append(w3calls, eth.EstimateGas(&msg, nil).Returns(&estimatedGasses[i]))

		// simulation will not track the state of the blockchain after approval
		if strings.Contains(strings.ToLower(strings.ReplaceAll(call.CxFunc, " ", "")), "approve") {
			includesApprove = true
		}
	}
	// Estimate gas of messages
	err := e.client.Call(w3calls...)
	_, ok := err.(w3.CallErrors)
	if !includesApprove && (err != nil && ok) {
		return CallEstimate{Value: *big.NewInt(0), Gas: 0, Success: false}, libcommon.StringError(err)
	}

	// Call(w3calls) should fill out estimatedGasses array
	totalGas := uint64(0)
	for _, gas := range estimatedGasses {
		totalGas += gas
	}
	return CallEstimate{Value: *totalValue, Gas: totalGas, Success: true}, nil
}

func (e executor) Initiate(calls []ContractCall) ([]string, *big.Int, error) {
	w3calls := []w3types.Caller{}
	hashes := make([]ethcommon.Hash, len(calls))
	totalValue := big.NewInt(0)
	for i, call := range calls {
		// TODO: Further optimize this to take in the calls array
		tx, err := e.generateTransactionRequest(call, uint64(i))
		if err != nil {
			return []string{}, nil, libcommon.StringError(err)
		}
		totalValue = totalValue.Add(totalValue, tx.Value())
		w3calls = append(w3calls, eth.SendTx(&tx).Returns(&hashes[i]))
	}

	// Call txs and retrieve hashes
	err := e.client.Call(w3calls...)
	callErrs, ok := err.(w3.CallErrors)
	if err != nil && ok {
		catErrs := ""
		for _, callErr := range callErrs {
			if callErr != nil {
				catErrs += callErr.Error() + " "
			}
		}
		return []string{}, nil, libcommon.StringError(errors.New(catErrs))
	}

	hashStrings := make([]string, len(hashes))
	for i := range hashes {
		hashStrings[i] = hashes[i].String()
	}
	return hashStrings, totalValue, nil
}

func (e executor) TxWait(txIds []string) (uint64, error) {
	totalGasUsed := uint64(0)
	for _, txId := range txIds {
		txHash := ethcommon.HexToHash(txId)
		receipt := types.Receipt{}
		pending := true
		for pending {
			pendingReceipt, err := e.geth.TransactionReceipt(context.Background(), txHash)
			// TransactionReceipt returns error "not found" while tx is pending
			if err != nil && err.Error() != "not found" {
				return totalGasUsed, libcommon.StringError(err)
			}
			if pendingReceipt != nil {
				receipt = *pendingReceipt
				pending = false
			}
			// TODO: Sleep for a few ms to keep the cpu cooler
		}
		if receipt.Status == 0 {
			return totalGasUsed, libcommon.StringError(errors.New("transaction failed"))
		}
		totalGasUsed += receipt.GasUsed
	}
	return totalGasUsed, nil
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

func (e executor) generateTransactionMessage(call ContractCall, incrementNonce uint64) (w3types.Message, error) {
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
		Nonce:     nonce + incrementNonce,
	}, nil
}

func (e executor) generateTransactionRequest(call ContractCall, incrementNonce uint64) (types.Transaction, error) {
	tx := types.Transaction{}

	msg, err := e.generateTransactionMessage(call, 0)
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
		Nonce:     msg.Nonce + incrementNonce,
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

func (e executor) GetEventData(txIds []string, eventSignature string) ([]types.Log, error) {
	events := []types.Log{}
	for _, txId := range txIds {
		receipt, err := e.geth.TransactionReceipt(context.Background(), ethcommon.HexToHash(txId))
		if err != nil {
			return []types.Log{}, libcommon.StringError(err)
		}

		event := crypto.Keccak256Hash([]byte(eventSignature))

		// Iterate through the logs to find the transfer event and extract the token ID.
		for _, log := range receipt.Logs {
			if log.Topics[0].Hex() == event.Hex() {
				events = append(events, *log)
			}
		}
	}
	return events, nil
}

// This can be used to check if the recipient of an event such as transfer matches our hot wallet address
func FilterEventData(logs []types.Log, indexes []int, hexValues []string) []types.Log {
	matches := []types.Log{}
	for _, log := range logs {
		for i, index := range indexes {
			// Event Address is checksummed, but Event Topics are not
			// compare RHS of topic with hexValues query
			expected := hexValues[i]
			if expected[:2] == "0x" {
				expected = expected[2:]
			}
			RHS := log.Topics[index].Hex()[len(log.Topics[index].Hex())-len(expected):]
			if strings.EqualFold(RHS, expected) {
				matches = append(matches, log)
			}
		}
	}
	return matches
}

func (e executor) GetTokenIds(txIds []string) ([]string, error) {
	logs, err := e.GetEventData(txIds, "Transfer(address,address,uint256)")
	if err != nil {
		return []string{}, libcommon.StringError(err)
	}
	tokenIds := []string{}
	for _, log := range logs {
		if len(log.Topics) != 4 {
			continue
		}
		tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes())
		tokenIds = append(tokenIds, tokenId.String())
	}
	return tokenIds, nil
}

func (e executor) GetTokenQuantities(txIds []string) ([]string, error) {
	logs, err := e.GetEventData(txIds, "Transfer(address,address,uint256)")
	if err != nil {
		return []string{}, libcommon.StringError(err)
	}
	quantities := []string{}
	for _, log := range logs {
		quantity := new(big.Int).SetBytes(log.Data)
		quantities = append(quantities, quantity.String())
	}
	return quantities, nil
}

func (e executor) ForwardNonFungibleTokens(txIds []string, recipient string) ([]string, []string, error) {
	eventData, err := e.GetEventData(txIds, "Transfer(address,address,uint256)")
	if err != nil {
		return []string{}, []string{}, libcommon.StringError(err)
	}
	hotWallet, err := e.getAccount()
	if err != nil {
		return []string{}, []string{}, libcommon.StringError(err)
	}
	// Filter events where recipient is our hot wallet
	toForward := FilterEventData(eventData, []int{2}, []string{hotWallet.String()})
	tokenIds := []string{}
	calls := []ContractCall{}
	for _, log := range toForward {
		tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes()).String()
		call := ContractCall{
			CxAddr: log.Address.String(),
			CxFunc: "transferFrom(address,address,uint256)",
			CxParams: []string{
				hotWallet.String(),
				recipient,
				tokenId,
			},
			CxReturn:   "",
			TxValue:    "0",
			TxGasLimit: "800000",
		}
		tokenIds = append(tokenIds, tokenId)
		calls = append(calls, call)
	}

	forwardTxIds := []string{}
	if len(calls) > 0 {
		forwardTxIds, _, err = e.Initiate(calls)
		if err != nil {
			return txIds, tokenIds, libcommon.StringError(err)
		}
	}

	return forwardTxIds, tokenIds, nil
}

func (e executor) ForwardTokens(txIds []string, recipient string) ([]string, []string, []string, error) {
	eventData, err := e.GetEventData(txIds, "Transfer(address,address,uint256)")
	if err != nil {
		return []string{}, []string{}, []string{}, libcommon.StringError(err)
	}
	hotWallet, err := e.getAccount()
	if err != nil {
		return []string{}, []string{}, []string{}, libcommon.StringError(err)
	}
	// Filter events where recipient is our hot wallet
	toForward := FilterEventData(eventData, []int{2}, []string{hotWallet.String()})
	tokens := []string{}
	quantities := []string{}
	calls := []ContractCall{}
	for _, log := range toForward {
		quantity := new(big.Int).SetBytes(log.Data).String()
		token := log.Address.String()
		call := ContractCall{
			CxAddr: token,
			CxFunc: "transfer(address,uint256)",
			CxParams: []string{
				recipient,
				quantity,
			},
			CxReturn:   "",
			TxValue:    "0",
			TxGasLimit: "800000",
		}
		tokens = append(tokens, token)
		quantities = append(quantities, quantity)
		calls = append(calls, call)
	}

	forwardTxIds := []string{}
	if len(calls) > 0 {
		forwardTxIds, _, err = e.Initiate(calls)
		if err != nil {
			return txIds, tokens, quantities, libcommon.StringError(err)
		}
	}

	return forwardTxIds, tokens, quantities, nil
}
