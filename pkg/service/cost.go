package service

import (
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type EstimationParams struct {
	ChainID    uint64  `json:"chainID"`
	CostETH    big.Int `json:"costETH"`
	UseBuffer  bool    `json:"useBuffer"`
	GasUsedWei uint64  `json:"gasUsedWei"`
	CostToken  big.Int `json:"costToken"`
	TokenName  string  `json:"tokenName"`
}

type OwlracleJSON struct {
	Timestamp string  `json:"timestamp"`
	LastBlock uint64  `json:"lastBlock"`
	AvgTime   float64 `json:"avgTime"`
	AvgTx     float64 `json:"avgTx"`
	AvgGas    float64 `json:"avgGas"`
	Speeds    []struct {
		Acceptance           float64 `json:"acceptance"`
		MaxFeePerGas         float64 `json:"maxFeePerGas"`
		MaxPriorityFeePerGas float64 `json:"maxPriorityFeePerGas"`
		BaseFee              float64 `json:"baseFee"`
		EstimatedFee         float64 `json:"estimatedFee"`
	} `json:"speeds"`
}

type Cost interface {
	EstimateTransaction(p EstimationParams) (model.Quote, error)
	New(repo repository.Cost) Cost
	QueryOwlracle(chainId uint64) (float64, error)
	LookupUSD(coin string, quantity float64) (float64, error)
}

type cost struct {
	repository repository.Cost // cached token and gas costs
}

func (c cost) New(repo repository.Cost) Cost {
	return &cost{
		repository: repo,
	}
}

func NewCost(repo repository.Cost) Cost {
	return &cost{
		repository: repo,
	}
}

func (c cost) QueryOwlracle(chainId uint64) (float64, error) {
	blockChain, err := model.ChainInfo(chainId)
	if err != nil {
		return 0, err
	}
	gwei, err := c.owlracle(blockChain.OwlracleName)
	if err != nil {
		return 0, err
	}
	return gwei, nil
}

func (c cost) EstimateTransaction(p EstimationParams) (model.Quote, error) {
	// Get Unix Timestamp and chain info
	timestamp := time.Now().Unix()
	blockChain, err := model.ChainInfo(p.ChainID)
	if err != nil {
		return model.Quote{}, err
	}

	// Query cost of native token in USD
	nativeCost, err := c.LookupUSD(blockChain.CoingeckoName, 1)
	if err != nil {
		return model.Quote{}, err
	}

	// Use it to convert transactioncost and apply buffer
	if p.UseBuffer {
		nativeCost *= 1.0 + common.NativeTokenBuffer(blockChain.ChainID)
	}
	costEth := common.WeiToEther(&p.CostETH)
	transactionCost := costEth * nativeCost

	// Query owlracle for gas
	ethGasFee, err := c.lookupGas(blockChain.OwlracleName)
	if err != nil {
		return model.Quote{}, err
	}

	// Convert it from gwei to eth to USD and apply buffer
	gasInUSD := ethGasFee * float64(p.GasUsedWei) * nativeCost / float64(1e9)
	if p.UseBuffer {
		gasInUSD *= 1.0 + common.GasBuffer(blockChain.ChainID)
	}

	// Query cost of token in USD if used and apply buffer
	costToken := common.WeiToEther(&p.CostToken)
	tokenCost, err := c.LookupUSD(p.TokenName, costToken)
	if err != nil {
		return model.Quote{}, err
	}
	if p.UseBuffer {
		tokenCost *= 1.0 + common.TokenBuffer(p.TokenName)
	}

	// Compute service fee
	upcharge := blockChain.StringFee
	serviceFee := (transactionCost + gasInUSD + tokenCost) * upcharge

	// Fill out CostEstimate and return
	return model.Quote{
		Timestamp:  timestamp,
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

func (c cost) LookupUSD(coin string, quantity float64) (float64, error) {
	// DB under construction
	res, err := c.coingeckoUSD(coin, 1)
	if err != nil {
		return 0, err
	}
	return res * quantity, nil
}

func (c cost) lookupGas(network string) (float64, error) {
	// DB under construction
	res, err := c.owlracle(network)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (c cost) coingeckoUSD(coin string, quantity float64) (float64, error) {
	requestURL := os.Getenv("COINGECKO_API_URL") + "simple/price?ids=" + coin + "&vs_currencies=usd"
	var res map[string]interface{}
	err := common.GetJson(requestURL, &res)
	if err != nil {
		return 0, err
	}
	prices, found := res[coin]
	if found {
		priceMap := prices.(map[string]interface{})
		usd, found := priceMap["usd"]
		if found {
			return usd.(float64), nil
		}
	}
	return 0, nil
}

func (c cost) owlracle(network string) (float64, error) {
	requestURL := os.Getenv("OWLRACLE_API_URL") +
		network +
		"/gas?apikey=" +
		os.Getenv("OWLRACLE_API_KEY") +
		"&accept=100"
	var res OwlracleJSON
	err := common.GetJson(requestURL, &res)
	if err != nil {
		return 0, err
	}
	if len(res.Speeds) > 0 {
		return res.Speeds[0].MaxFeePerGas, nil
	}
	return 0, errors.New("owlracle: malformed response")
}
