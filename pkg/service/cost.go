package service

import (
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type EstimationParams struct {
	ChainID     uint32  `json:"chainID"`
	CostETH     float32 `json:"costETH"`
	UseBuffer   bool    `json:"useBuffer"`
	GasUsedGwei uint32  `json:"gasUsedGwei"`
	CostToken   float32 `json:"costToken"`
	TokenName   string  `json:"tokenName"`
}

type CostEstimate struct {
	Timestamp  int64   `json:"timeStamp"`
	BaseUSD    float32 `json:"baseUSD"`
	GasUSD     float32 `json:"gasUSD"`
	TokenUSD   float32 `json:"tokenUSD"`
	ServiceUSD float32 `json:"serviceUSD"`
	TotalUSD   float32 `json:"totalUSD"`
}

type SignedQuote struct {
	Estimate  CostEstimate
	Signature string `json:"signature"`
}

// type coingeckoJSON struct {
// 	prices []string {
// 		currency [] string : float32
// 	}
// }

type Cost interface {
	EstimateTransaction(p EstimationParams) (CostEstimate, error)
	New(repo repository.Cost) Cost
}

type cost struct {
	repository repository.Cost // cached token and gas costs
}

func (c cost) New(repo repository.Cost) Cost {
	return &cost{repository: repo}
}

func NewCost(repo repository.Cost) Cost {
	return &cost{repository: repo}
}

func (c cost) EstimateTransaction(p EstimationParams) (CostEstimate, error) {
	// Get Unix Timestamp
	date := time.Now().Unix()
	blockChain, err := model.ChainInfo(p.ChainID)
	if err != nil {
		return CostEstimate{}, err
	}
	// Query cost of native token in USD
	nativeCost := c.getUSDFromDB(blockChain.CoingeckoName)
	// Use it to convert transactioncost and apply buffer
	if p.UseBuffer {
		nativeCost *= 1.0 + common.NativeTokenBuffer(blockChain.ChainID)
	}
	transactionCost := p.CostETH * nativeCost
	// Query owlracle for gas
	ethGasFee := c.getGasFromDB(blockChain.OwlracleName)
	// Convert it from gwei to eth and apply buffer
	gasInUSD := ethGasFee * float32(p.GasUsedGwei) * nativeCost / 1e9
	if p.UseBuffer {
		gasInUSD *= 1.0 + common.GasBuffer(blockChain.ChainID)
	}
	// Query cost of token in USD if used and apply buffer
	tokenCost := c.getUSDFromDB(p.TokenName) * p.CostToken
	if p.UseBuffer {
		tokenCost *= 1.0 + common.TokenBuffer(p.TokenName)
	}
	// Compute service fee
	upcharge := blockChain.StringFee
	serviceFee := (transactionCost + gasInUSD + tokenCost) * upcharge
	// Fill out CostEstimate and return
	return CostEstimate{
		Timestamp:  date,
		BaseUSD:    transactionCost,
		GasUSD:     gasInUSD,
		TokenUSD:   tokenCost,
		ServiceUSD: serviceFee,
		TotalUSD:   transactionCost + gasInUSD + tokenCost + serviceFee,
	}, nil
}

func (c cost) getExternalAPICallInterval(rateLimitPerMinute float32, uniqueEntries uint32) float32 {
	return (float32(uniqueEntries*60000) / rateLimitPerMinute)
}

func (c cost) getUSDFromDB(coin string) float32 {
	return 0
}

func (c cost) getGasFromDB(network string) float32 {
	return 0
}

func coingeckoUSD(coin string, quantity float32) (float32, error) {
	// requestURL := os.Getenv("COINGECKO_API_URL") + "simple/price?ids=" + coin + "&vs_currencies=usd"
	// response, err := http.Get(requestURL)
	// if err != nil {
	// 	return 0, err
	// }
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	return 0, err
	// }
	// body
	return 0, nil
}

func owlracle(network string) float32 {
	return 0
}
