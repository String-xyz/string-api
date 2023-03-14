package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	commonlib "github.com/String-xyz/go-lib/common"
	"github.com/ethereum/go-ethereum/accounts"
	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
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
		return ethcomm.Address{}, commonlib.StringError(err)
	}
	return crypto.PubkeyToAddress(*recovered), nil
}

func BigNumberToFloat(bigNumber string, decimals uint64) (floatReturn float64, err error) {
	floatReturn, err = strconv.ParseFloat(bigNumber, 64)
	if err != nil {
		log.Err(err).Msg("Failed to convert bigNumber to float")
		err = commonlib.StringError(err)
		return
	}
	floatReturn = floatReturn * math.Pow(10, -float64(decimals))
	return
}

func GetBaseURL() string {
	return os.Getenv("BASE_URL")
}

func FloatToUSDString(amount float64) string {
	return fmt.Sprintf("USD $%.2f", math.Round(amount*100)/100)
}

func BetterStringify(jsonBody any) (betterString string, err error) {
	bodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		log.Err(err).Interface("body", jsonBody).Msg("Could not encode to bytes")
		return betterString, commonlib.StringError(err)
	}

	bodyReader := bytes.NewReader(bodyBytes)

	betterBytes, err := io.ReadAll(bodyReader)
	betterString = string(betterBytes)
	if err != nil {
		return betterString, commonlib.StringError(err)
	}

	return
}

func SliceContains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
