package service

import (
	"fmt"
	"math/big"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/lmittmann/w3"
	"github.com/pkg/errors"
)

type SwapQuoteData struct {
	yield   *big.Int
	gas     *big.Int
	chainId int
}

// Hackathon
func nativeTokenOracleName(chainId int) string {
	switch chainId {
	case 1:
		return "ethereum"
	case 43114:
		return "avalanche-2"
	}
	return "unknown"
}

func SwapQuote(from string, to string, chainId int, amount *big.Int) (SwapQuoteData, error) {
	quoteData := SwapQuoteData{}
	request :=
		"https://api.1inch.io/v5.0/" +
			fmt.Sprint(chainId) +
			"/quote?fromTokenAddress=" +
			from + "&toTokenAddress=" +
			to + "&amount=" + fmt.Sprint(amount.String())
	var res map[string]interface{}

	err := common.GetJsonGeneric(request, &res)
	if err != nil {
		return quoteData, common.StringError(err)
	}

	toAmount, valid := res["toTokenAmount"]
	if !valid {
		return quoteData, common.StringError(err)
	}
	gasAmount, valid := res["estimatedGas"]
	if !valid {
		return quoteData, common.StringError(err)
	}
	quoteData.yield = w3.I(toAmount.(string))
	quoteData.gas = w3.I(fmt.Sprint(gasAmount.(float64)))
	quoteData.chainId = chainId
	return quoteData, nil
}

func SwapQuoteToUSD(quoteData SwapQuoteData) (float64, error) {
	// Hackathon
	c := NewCost(nil)

	// TODO: Use Cost.EstimateTransaction
	oracleName := nativeTokenOracleName(quoteData.chainId)
	nativeValueUSD, err := c.CoingeckoUSD(oracleName, 1)
	if err != nil {
		return 0, common.StringError(err)
	}
	yieldMinusGas := quoteData.yield.Sub(quoteData.yield, quoteData.gas) // hackathon
	ethValue := common.WeiToEther(yieldMinusGas)
	return ethValue * nativeValueUSD, nil
}

func SwapQuoteUSD(from string, to string, chainId int, amount *big.Int) (float64, error) {
	quoteData, err := SwapQuote(from, to, chainId, amount)
	if err != nil {
		return 0, common.StringError(err)
	}
	usd, err := SwapQuoteToUSD(quoteData)
	if err != nil {
		return 0, common.StringError(err)
	}
	return usd, nil
}

func SwapQuoteCrossChain(from string, fromChain int, to string, toChain int, fromAmount *big.Int) (float64, error) {
	fromValueUSD, err := SwapQuoteUSD(from, "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", fromChain, fromAmount)
	if err != nil {
		return 0, common.StringError(err)
	}
	c := NewCost(nil) // Hackathon
	destinationNativeTokenCostUSD, err := c.CoingeckoUSD(nativeTokenOracleName(toChain), 1)
	if err != nil {
		return 0, common.StringError(err)
	}
	destinationNativeTokenAmount := fromValueUSD / destinationNativeTokenCostUSD
	destinationNativeTokenWei := w3.I(fmt.Sprintf("%f ether", destinationNativeTokenAmount)) // Hackathon
	destinationQuoteData, err := SwapQuote("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", to, toChain, destinationNativeTokenWei)
	if err != nil {
		return 0, common.StringError(err)
	}
	yieldMinusGas := destinationQuoteData.yield.Sub(destinationQuoteData.yield, destinationQuoteData.gas) // hackathon
	destinationTokenAmount := common.WeiToEther(yieldMinusGas)
	return destinationTokenAmount, nil
}

func OffRamp(from string, fromChain int, fromAmount *big.Int) {
	// Receive end users tokens

	// Trade them for native token

	// Send USD to the end users CC
}

func Swap(from string, to string, chainId int, userAddr string, amount *big.Int) error {
	var spender map[string]interface{} // Address of 1inch smart contract

	getSpenderAddr := "https://api.1inch.io/v5.0/" + fmt.Sprint(chainId) + "/approve/spender"
	err := common.GetJsonGeneric(getSpenderAddr, &spender)
	if err != nil {
		return common.StringError(err)
	}
	addr, ok := spender["address"].(string)
	if !ok {
		return common.StringError(errors.New("Failed to unmarshal 1inch spender addr"))
	}
	addr = common.SanitizeChecksum(addr) // 1inch provides non-checksummed address

	e := NewExecutor()
	e.Initialize("https://api.avax.network/ext/bc/C/rpc") // Hackathon
	call := ContractCall{
		CxAddr:     from,
		CxFunc:     "approve(address,uint256)",
		CxReturn:   "bool",
		CxParams:   []string{addr, amount.String()},
		TxValue:    "0",
		TxGasLimit: "8000000",
	}
	approveTx, approveGas, err := e.Initiate(call)
	if err != nil {
		return common.StringError(err)
	}
	_, err = e.TxWait(approveTx)
	if err != nil {
		return common.StringError(err)
	}
	fmt.Printf("\nUsed %+v gas to approve swap %+v", common.WeiToEther(approveGas), approveTx)

	request := "https://api.1inch.io/v5.0/" +
		fmt.Sprint(chainId) +
		"/swap?fromTokenAddress=" + from +
		"&toTokenAddress=" + to +
		"&amount=" + amount.String() +
		"&fromAddress=" + userAddr +
		"&slippage=1"

	var res map[string]interface{}

	err = common.GetJsonGeneric(request, &res)
	if err != nil {
		return common.StringError(err)
	}

	tx, ok := res["tx"].(map[string]interface{})
	if !ok {
		return common.StringError(err)
	}

	// HACKATHON GO GO GO
	data := tx["data"].(string)
	fromm := tx["from"].(string)
	gas := tx["gas"].(float64)
	gasPrice := tx["gasPrice"].(string)
	too := tx["to"].(string)
	value := tx["value"].(string)

	fromm = common.SanitizeChecksum(fromm)
	too = common.SanitizeChecksum(too)
	swap := EncodedContractCall{
		Data:     data,
		From:     fromm,
		Gas:      fmt.Sprint(gas),
		GasPrice: *w3.I(gasPrice),
		To:       too,
		Value:    value,
	}
	txId, swapGas, err := e.InitiateEncoded(swap)
	fmt.Printf("\n\nINITIATED SWAP %+v", txId)
	if err != nil {
		return common.StringError(err)
	}
	_, err = e.TxWait(txId)
	if err != nil {
		return common.StringError(err)
	}

	fmt.Printf("\n\nSWAPPED TX = %+v using %+v gas", txId, swapGas.String())
	return nil
}
