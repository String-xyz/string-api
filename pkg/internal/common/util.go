package common

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"
	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func ToSha256(v string) string {
	bs := sha256.Sum256([]byte(v))
	return hex.EncodeToString(bs[:])
}

func RecoverAddress(message string, signature string) (ethcomm.Address, error) {
	sig := hexutil.MustDecode(signature)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27
	}
	msg := accounts.TextHash([]byte(message))
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return ethcomm.Address{}, StringError(err)
	}
	return crypto.PubkeyToAddress(*recovered), nil
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// keysAndValues is only being used for optional updates
// do not use it for insert or select
func KeysAndValues(item interface{}) ([]string, map[string]interface{}) {
	tag := "db"
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	keyNames := make([]string, 0, v.NumField())
	keyValues := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := reflectValue.Field(i).Interface()
		if !isNil(field) {
			t := v.Field(i).Tag.Get(tag) + "=:" + v.Field(i).Tag.Get(tag)
			keyNames = append(keyNames, t)
			keyValues[v.Field(i).Tag.Get(tag)] = field
		}
	}

	return keyNames, keyValues
}
