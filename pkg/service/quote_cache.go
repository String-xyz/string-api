package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/pkg/errors"
)

type QuoteCache interface {
	CheckUpdateCachedTransactionRequest(request model.TransactionRequest, desiredInterval int64) (recalculate bool, callEstimate CallEstimate, err error)
	PutCachedTransactionRequest(request model.TransactionRequest, data CallEstimate) error
}

type quoteCache struct {
	redis database.RedisStore
}

func NewQuoteCache(redis database.RedisStore) QuoteCache {
	return &quoteCache{redis}
}

type callEstimateCache struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value" db:"value"` // We use string here to store in db
	Gas       uint64 `json:"gas" db:"gas"`
	Success   bool   `json:"success" db:"success"`
}

func (q quoteCache) CheckUpdateCachedTransactionRequest(request model.TransactionRequest, desiredInterval int64) (recalculate bool, callEstimate CallEstimate, err error) {
	cacheObject, err := store.GetObjectFromCache[callEstimateCache](q.redis, tokenizeTransactionRequest(sanitizeTransactionRequest(request)))
	if cacheObject.Timestamp == 0 || (err == nil && time.Now().Unix()-cacheObject.Timestamp > desiredInterval) {
		return true, CallEstimate{}, nil
	} else if err != nil {
		return false, CallEstimate{}, libcommon.StringError(err)
	} else {
		// Construct CallEstimate from cacheObject
		value := new(big.Int)
		value, ok := value.SetString(cacheObject.Value, 10)
		if !ok {
			return false, CallEstimate{}, libcommon.StringError(errors.New("Failed to parse value from cache"))
		}
		return false, CallEstimate{Value: *value, Gas: cacheObject.Gas, Success: cacheObject.Success}, nil
	}
}

func (q quoteCache) PutCachedTransactionRequest(request model.TransactionRequest, data CallEstimate) error {
	cacheObject := callEstimateCache{
		Timestamp: time.Now().Unix(),
		Value:     data.Value.String(),
		Gas:       data.Gas,
		Success:   data.Success,
	}
	err := store.PutObjectInCache(q.redis, tokenizeTransactionRequest(sanitizeTransactionRequest(request)), cacheObject)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

func sanitizeTransactionRequest(request model.TransactionRequest) model.TransactionRequest {
	// Structs are pointers
	sanitized := model.TransactionRequest{
		UserAddress: request.UserAddress,
		ChainId:     request.ChainId,
		CxAddr:      request.CxAddr,
		CxFunc:      request.CxFunc,
		CxReturn:    request.CxReturn,
		CxParams:    append([]string{}, request.CxParams...), // So are arrays
		TxValue:     request.TxValue,                         // TODO: Maybe omit this and take it in from the endpoint before converting to USD
		TxGasLimit:  request.TxGasLimit,
	}
	// Treat the users address as a wildcard
	for i, param := range sanitized.CxParams {
		if param == sanitized.UserAddress {
			sanitized.CxParams[i] = "*"
		}
	}
	sanitized.UserAddress = "*"
	return sanitized
}

func tokenizeTransactionRequest(request model.TransactionRequest) string {
	bytes, err := json.Marshal(request)
	if err != nil {
		return ""
	}
	hash := sha1.New()
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}
