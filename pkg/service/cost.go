package service

import (
	"math/big"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/pkg/errors"
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

type CostCache struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Cost interface {
	EstimateTransaction(p EstimationParams, chain Chain) (model.Quote, error)
	New(redis store.RedisStore) Cost
	LookupUSD(coin string, quantity float64) (float64, error)
}

type cost struct {
	redis store.RedisStore // cached token and gas costs
}

func (c cost) New(redis store.RedisStore) Cost {
	return &cost{
		redis: redis,
	}
}

func NewCost(redis store.RedisStore) Cost {
	return &cost{
		redis: redis,
	}
}

func (c cost) EstimateTransaction(p EstimationParams, chain Chain) (model.Quote, error) {
	// Get Unix Timestamp and chain info
	timestamp := time.Now().Unix()

	// Query cost of native token in USD
	nativeCost, err := c.LookupUSD(chain.CoingeckoName, 1)
	if err != nil {
		return model.Quote{}, common.StringError(err)
	}

	// Use it to convert transactioncost and apply buffer
	if p.UseBuffer {
		nativeCost *= 1.0 + common.NativeTokenBuffer(chain.ChainID)
	}
	costEth := common.WeiToEther(&p.CostETH)
	transactionCost := costEth * nativeCost

	// Query owlracle for gas
	ethGasFee, err := c.lookupGas(chain.OwlracleName)
	if err != nil {
		return model.Quote{}, common.StringError(err)
	}

	// Convert it from gwei to eth to USD and apply buffer
	gasInUSD := ethGasFee * float64(p.GasUsedWei) * nativeCost / float64(1e9)
	if p.UseBuffer {
		gasInUSD *= 1.0 + common.GasBuffer(chain.ChainID)
	}

	// Query cost of token in USD if used and apply buffer
	costToken := common.WeiToEther(&p.CostToken)
	tokenCost, err := c.LookupUSD(p.TokenName, costToken)
	if err != nil {
		return model.Quote{}, common.StringError(err)
	}
	if p.UseBuffer {
		tokenCost *= 1.0 + common.TokenBuffer(p.TokenName)
	}

	// Compute service fee
	upcharge := chain.StringFee
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

func (c cost) getExternalAPICallInterval(rateLimitPerMinute float32, uniqueEntries uint32) int64 {
	return int64(float32(uniqueEntries*60000) / rateLimitPerMinute)
}

func (c cost) LookupUSD(coin string, quantity float64) (float64, error) {
	cacheName := "usd_value_" + coin
	cacheObject, err := GetObjectFromCache[CostCache](c.redis, cacheName)
	if err != nil && errors.Cause(err).Error() != "redis: nil" {
		return 0.0, common.StringError(err)
	}
	if cacheObject == (CostCache{}) || (err == nil && time.Now().Unix()-cacheObject.Timestamp > c.getExternalAPICallInterval(10, 6)) {
		cacheObject.Timestamp = time.Now().Unix()
		cacheObject.Value, err = c.coingeckoUSD(coin, 1)
		if err != nil {
			return 0, common.StringError(err)
		}
		err = PutObjectInCache(c.redis, cacheName, cacheObject)
		if err != nil {
			return 0.0, common.StringError(err)
		}
	}

	return cacheObject.Value * quantity, nil
}

func (c cost) lookupGas(network string) (float64, error) {
	cacheName := "gas_price_" + network
	cacheObject, err := GetObjectFromCache[CostCache](c.redis, cacheName)
	if err != nil {
		return 0.0, common.StringError(err)
	}
	if cacheObject == (CostCache{}) || time.Now().Unix()-cacheObject.Timestamp > c.getExternalAPICallInterval(1.6, 6) {
		cacheObject.Timestamp = time.Now().Unix()
		cacheObject.Value, err = c.owlracle(network)
		if err != nil {
			return 0, common.StringError(err)
		}
		err = PutObjectInCache(c.redis, cacheName, cacheObject)
		if err != nil {
			return 0.0, common.StringError(err)
		}
	}

	return cacheObject.Value, nil
}

func (c cost) coingeckoUSD(coin string, quantity float64) (float64, error) {
	requestURL := os.Getenv("COINGECKO_API_URL") + "simple/price?ids=" + coin + "&vs_currencies=usd"
	var res map[string]interface{}
	err := common.GetJson(requestURL, &res)
	if err != nil {
		return 0, common.StringError(err)
	}
	prices, found := res[coin]
	if found {
		priceMap := prices.(map[string]interface{})
		usd, found := priceMap["usd"]
		if found {
			return usd.(float64), nil
		}
	}
	// return 0, common.StringError(errors.New("Price not found for " + coin))
	// fmt.Printf("\n\nPRICE LOOKUP %+v", coin)
	// TODO: this is getting hit somewhere, figure out why
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
		return 0, common.StringError(err)
	}
	if len(res.Speeds) > 0 {
		return res.Speeds[0].MaxFeePerGas, nil
	}
	return 0, errors.New("owlracle: malformed response")
}
