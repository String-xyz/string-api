package common

import (
	"context"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/String-xyz/go-lib/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/lmittmann/w3"
	"golang.org/x/crypto/sha3"
)

func ParseEncoding(function *w3.Func, signature string, params []string) ([]byte, error) {
	signatureArgs := strings.Split(strings.Split(strings.Split(signature, "(")[1], ")")[0], ",")
	if len(signatureArgs) != len(params) {
		return nil, common.StringError(errors.New("executor parseParams: mismatched arguments"))
	}
	args := []interface{}{}
	for i, s := range signatureArgs {
		switch s {
		case "address":
			args = append(args, w3.A(params[i]))
		case "bool":
			args = append(args, params[i] == "true" || params[i] == "TRUE")
		case "string":
			args = append(args, params[i])
		case "bytes":
			args = append(args, w3.B(params[i]))
		case "uint8":
			v, err := strconv.ParseUint(params[i], 0, 8)
			if err != nil {
				return nil, common.StringError(err)
			}
			args = append(args, v)
		case "uint32":
			v, err := strconv.ParseUint(params[i], 0, 32)
			if err != nil {
				return nil, common.StringError(err)
			}
			args = append(args, v)
		case "uint256":
			args = append(args, w3.I(params[i]))
		case "int8":
			v, err := strconv.ParseInt(params[i], 0, 8)
			if err != nil {
				return nil, common.StringError(err)
			}
			args = append(args, v)
		case "int32":
			v, err := strconv.ParseInt(params[i], 0, 32)
			if err != nil {
				return nil, common.StringError(err)
			}
			args = append(args, v)
		case "int256":
			args = append(args, w3.I(params[i]))
		default:
			return nil, common.StringError(errors.New("executor: parseParams: unsupported type"))
		}
	}
	result, err := function.EncodeArgs(args...)
	if err != nil {
		return nil, common.StringError(err)
	}
	return result, nil
}

func WeiToEther(wei *big.Int) float64 {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)
	fWei.SetMode(big.ToNearestEven)
	ethBig := f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
	eth64, _ := ethBig.Float64() // OK to reduce precision?
	return eth64
}

// TODO: Eventually make sure we support smart contract wallets
func IsWallet(addr string) bool {
	RPC := "https://rpc.ankr.com/eth" // temporarily just use ETH mainnet
	geth, _ := ethclient.Dial(RPC)

	if !validAddress(addr) {
		return false
	}
	addr = SanitizeChecksum(addr) // Copy correct checksum, although endpoint handlers are doing this already

	address := ethcommon.HexToAddress(addr)
	bytecode, err := geth.CodeAt(context.Background(), address, nil)
	if err != nil {
		return false
	}
	isContract := len(bytecode) > 0
	return !isContract
}

func validChecksum(addr string) bool {
	valid := SanitizeChecksum(addr)
	return addr == valid
}

func SanitizeChecksum(addr string) string {
	lowerCase := strings.ToLower(addr)[2:]
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(lowerCase))
	hashBytes := hash.Sum(nil)

	valid := "0x"
	for i, b := range lowerCase {
		c := string(b)
		if b < '0' || b > '9' {
			if hashBytes[i/2]&byte(128-i%2*120) != 0 {
				c = string(b - 32)
			}
		}
		valid += c
	}
	return valid
}

func validAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}
