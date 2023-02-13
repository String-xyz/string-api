package service

import (
	"fmt"
	"testing"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/joho/godotenv"
	"github.com/lmittmann/w3"
	"github.com/stretchr/testify/assert"
)

func TestGetSwapQuote(t *testing.T) {
	quoteData, err := SwapQuote(
		"0xC38f41A296A4493Ff429F1238e030924A1542e50", // The token we have (SNOB)
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // The token we want to exchange it for (Native Token)
		43114,              // The chain we are on
		w3.I("1000 ether")) // The amount of the token we want to swap (1000 SNOB)
	assert.NoError(t, err)
	fmt.Printf("1000 SNOB yields us %+v AVAX minus %+v gas", quoteData.yield.String(), quoteData.gas.String())
}

func TestGetSwapQuoteUSD(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)
	usd, err := SwapQuoteUSD(
		"0xC38f41A296A4493Ff429F1238e030924A1542e50", // The token we have (SNOB)
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // The token we want to exchange it for (Native Token)
		43114,              // The chain we are on
		w3.I("1000 ether")) // The amount of the token we want to swap (1000 SNOB)
	assert.NoError(t, err)
	fmt.Printf("1000 SNOB yields us = %+v worth of AVAX", common.FloatToUSDString(usd))
	assert.NotEqual(t, 0, usd)
}

func TestGetCrossChainQuote(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	res, err := SwapQuoteCrossChain(
		"0xC38f41A296A4493Ff429F1238e030924A1542e50", // Token we have (SNOB)
		43114, // Chain we are starting from (Avalanche)
		"0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce", // Token we want (SHIB)
		1,                  // Chain we are swapping into
		w3.I("1000 ether")) // Amount of starting token we have (1000 SNOB)
	assert.NoError(t, err)

	fmt.Printf("1000 SNOB yields us = %.2f SHIBA INU", res)
	assert.NotEqual(t, 0, res)
}

func TestGetSwapPayload(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	err = Swap(
		"0xC38f41A296A4493Ff429F1238e030924A1542e50", // The token we have (SNOB)
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // The token we want to exchange it for (Native Token)
		43114, // The chain we are on
		"0x44A4b9E2A69d86BA382a511f845CbF2E31286770", // The wallet address which will initiate the swap
		w3.I("10 ether")) // The amount of the token we want to swap (10 SNOB)
	assert.NoError(t, err)
}
