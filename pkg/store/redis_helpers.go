package store

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/pkg/errors"
)

func GetObjectFromCache[T any](redis RedisStore, key string) (T, error) {
	var result *T = new(T)
	bytes, err := redis.Get(key)
	if err != nil && errors.Cause(err).Error() == "redis: nil" && len(bytes) == 0 {
		return *result, nil // object doesn't exist yet, create it down the stack
	} else if err != nil {
		// Work around the way that redis go api scopes error
		return *result, common.StringError(errors.New(err.Error()))
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return *result, common.StringError(err)
	}
	return *result, nil
}

func PutObjectInCache(redis RedisStore, key string, object any, optionalTimeout ...time.Duration) error {
	// Safeguard against missing tags
	val := reflect.ValueOf(object)
	for i := 0; i < val.Type().NumField(); i++ {
		if val.Type().Field(i).Tag.Get("json") == "" {
			return common.StringError(errors.New("object missing json tags"))
		}
	}

	var timeout time.Duration = 0
	if len(optionalTimeout) > 0 {
		timeout = optionalTimeout[0]
	}

	bytes, err := json.Marshal(object)
	if err != nil {
		return common.StringError(err)
	}

	err = redis.Set(key, bytes, timeout)
	if err != nil {
		// Work around the way that redis go API scopes error
		return common.StringError(errors.New(err.Error()))
	}
	return nil
}
