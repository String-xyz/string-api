package service

import (
	"math"
	"math/big"
	"os"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/pkg/errors"
)

type EstimationParams struct {
	ChainId    uint64  `json:"chainId"`
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
	LookupUSD(coin string, quantity float64) (float64, error)
}

type cost struct {
	redis database.RedisStore // cached token and gas costs
}

func NewCost(redis database.RedisStore) Cost {
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
		return model.Quote{}, libcommon.StringError(err)
	}

	// Use it to convert transactioncost and apply buffer
	if p.UseBuffer {
		nativeCost *= 1.0 + common.NativeTokenBuffer(chain.ChainId)
	}
	costEth := common.WeiToEther(&p.CostETH)
	// transactionCost is for native token transaction cost (tx_value)
	transactionCost := costEth * nativeCost

	// Query owlracle for gas
	ethGasFee, err := c.lookupGas(chain.OwlracleName)
	if err != nil {
		return model.Quote{}, libcommon.StringError(err)
	}

	// Convert it from gwei to eth to USD and apply buffer
	gasInUSD := ethGasFee * float64(p.GasUsedWei) * nativeCost / float64(1e9)
	if p.UseBuffer {
		gasInUSD *= 1.0 + common.GasBuffer(chain.ChainId)
	}

	// Query cost of token in USD if used and apply buffer
	costToken := common.WeiToEther(&p.CostToken)
	// tokenCost in contract call ERC-20 token costs
	// Also for buying tokens directly
	tokenCost, err := c.LookupUSD(p.TokenName, costToken)
	if err != nil {
		return model.Quote{}, libcommon.StringError(err)
	}
	if p.UseBuffer {
		tokenCost *= 1.0 + common.TokenBuffer(p.TokenName)
	}

	// Compute service fee
	upcharge := chain.StringFee
	baseCheckoutFee := 0.3
	serviceFee := (transactionCost+gasInUSD+tokenCost)*upcharge + baseCheckoutFee

	// floor
	if transactionCost < 0.01 {
		transactionCost = 0.01
	}
	if gasInUSD < 0.01 {
		gasInUSD = 0.01
	}

	// Round up to nearest cent
	transactionCost = centCeiling(transactionCost)
	gasInUSD = centCeiling(gasInUSD)
	tokenCost = centCeiling(tokenCost)
	serviceFee = centCeiling(serviceFee)

	// sum total
	totalUSD := transactionCost + gasInUSD + tokenCost + serviceFee

	// Round that up as well to account for any floating imprecision
	totalUSD = centCeiling(totalUSD)

	// Fill out CostEstimate and return
	return model.Quote{
		Timestamp:  timestamp,
		BaseUSD:    transactionCost,
		GasUSD:     gasInUSD,
		TokenUSD:   tokenCost,
		ServiceUSD: serviceFee,
		TotalUSD:   totalUSD,
	}, nil
}

func centCeiling(value float64) float64 {
	return math.Ceil(value*100) / 100
}

func (c cost) getExternalAPICallInterval(rateLimitPerMinute float64, uniqueEntries uint32) int64 {
	return int64(float64(60*rateLimitPerMinute) / rateLimitPerMinute)
}

func (c cost) LookupUSD(coin string, quantity float64) (float64, error) {
	cacheName := "usd_value_" + coin
	cacheObject, err := store.GetObjectFromCache[CostCache](c.redis, cacheName)
	if err != nil && serror.IsError(err, serror.NOT_FOUND) {
		return 0.0, libcommon.StringError(err)
	}
	if cacheObject == (CostCache{}) || (err == nil && time.Now().Unix()-cacheObject.Timestamp > c.getExternalAPICallInterval(10, 6)) {
		cacheObject.Timestamp = time.Now().Unix()
		// If coingecko is down, use coincap to get the price
		var empty interface{}
		err := common.GetJson(os.Getenv("COINGECKO_API_URL")+"ping", &empty)
		if err == nil {
			cacheObject.Value, err = c.coingeckoUSD(coin)
			if err != nil {
				return 0, libcommon.StringError(err)
			}
		} else {
			cacheObject.Value, err = c.coincapUSD(coin)
			if err != nil {
				return 0, libcommon.StringError(err)
			}
		}
		err = store.PutObjectInCache(c.redis, cacheName, cacheObject)
		if err != nil {
			return 0, libcommon.StringError(err)
		}
	}

	return cacheObject.Value * quantity, nil
}

func (c cost) lookupGas(network string) (float64, error) {
	cacheName := "gas_price_" + network
	cacheObject, err := store.GetObjectFromCache[CostCache](c.redis, cacheName)
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	if cacheObject == (CostCache{}) || time.Now().Unix()-cacheObject.Timestamp > c.getExternalAPICallInterval(1.6, 6) {
		cacheObject.Timestamp = time.Now().Unix()
		cacheObject.Value, err = c.owlracle(network)
		if err != nil {
			return 0, libcommon.StringError(err)
		}
		err = store.PutObjectInCache(c.redis, cacheName, cacheObject)
		if err != nil {
			return 0, libcommon.StringError(err)
		}
	}

	return cacheObject.Value, nil
}

func (c cost) coingeckoUSD(coin string) (float64, error) {
	requestURL := os.Getenv("COINGECKO_API_URL") + "simple/price?ids=" + coin + "&vs_currencies=usd"
	var res map[string]interface{}
	err := common.GetJsonGeneric(requestURL, &res)
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	prices, found := res[coin]
	if found {
		priceMap := prices.(map[string]interface{})
		usd, found := priceMap["usd"]
		if found {
			return usd.(float64), nil
		}
	}
	// return 0, libcommon.StringError(errors.New("Price not found for " + coin))
	// fmt.Printf("\n\nPRICE LOOKUP %+v", coin)
	// TODO: this is getting hit somewhere, figure out why
	return 0, nil
}

func (c cost) coincapUSD(coin string) (float64, error) {
	requestURL := os.Getenv("COINCAP_API_URL") + "assets?search=" + coin
	body := make(map[string]interface{})
	err := common.GetJsonGeneric(requestURL, &body)
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	res, found := body["data"].([]interface{})
	if found && len(res) > 0 {
		price, found := res[0].(map[string]interface{})["priceUsd"]
		if found {
			return price.(float64), nil
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
	err := common.GetJsonGeneric(requestURL, &res)
	if err != nil {
		return 0, libcommon.StringError(err)
	}
	if len(res.Speeds) > 0 {
		return res.Speeds[0].MaxFeePerGas, nil
	}
	return 0, errors.New("owlracle: malformed response")
}
