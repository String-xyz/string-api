package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"

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

func BigNumberToFloat(bigNumber string, decimals uint64) (floatReturn float64, err error) {
	floatReturn, err = strconv.ParseFloat(bigNumber, 64)
	if err != nil {
		log.Printf("Failed to convert bigNumber to float: %s", err)
		err = StringError(err)
		return
	}
	floatReturn = floatReturn * math.Pow(10, -float64(decimals))
	return
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

func GetBaseURL() string {
	return os.Getenv("BASE_URL")
}

func FloatToUSDString(amount float64) string {
	return fmt.Sprintf("USD $%.2f", math.Round(amount*100)/100)
}

func IsLocalEnv() bool {
	return os.Getenv("ENV") == "local"
}

func BetterStringify(jsonBody any) (betterString string, err error) {
	bodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Printf("Could not encode %+v to bytes: %s", jsonBody, err)
		return betterString, StringError(err)
	}

	bodyReader := bytes.NewReader(bodyBytes)

	betterBytes, err := io.ReadAll(bodyReader)
	betterString = string(betterBytes)
	if err != nil {
		return betterString, StringError(err)
	}

	return
}
