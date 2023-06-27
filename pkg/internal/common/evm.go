package common

import (
	"context"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/lmittmann/w3"
	"golang.org/x/crypto/sha3"
)

func addressArray(args []string) (set []ethcommon.Address) {
	for _, a := range args {
		set = append(set, w3.A(a))
	}
	return set
}

func boolArray(args []string) (set []bool) {
	for _, b := range args {
		set = append(set, b == "true" || b == "TRUE")
	}
	return set
}

func stringArray(args []string) (set []string) {
	return args
}

func bytesArray(args []string) (set [][]byte) {
	for _, b := range args {
		set = append(set, w3.B(b))
	}
	return set
}

func int8Array(args []string) (set []int8) {
	for _, i := range args {
		v, err := strconv.ParseInt(i, 0, 8)
		if err != nil {
			panic(err)
		}
		set = append(set, int8(v))
	}
	return set
}

func uint8Array(args []string) (set []uint8) {
	for _, u := range args {
		v, err := strconv.ParseUint(u, 0, 8)
		if err != nil {
			panic(err)
		}
		set = append(set, uint8(v))
	}
	return set
}

func uint16Array(args []string) (set []uint16) {
	for _, u := range args {
		v, err := strconv.ParseUint(u, 0, 16)
		if err != nil {
			panic(err)
		}
		set = append(set, uint16(v))
	}
	return set
}

func uint32Array(args []string) (set []uint32) {
	for _, u := range args {
		v, err := strconv.ParseUint(u, 0, 32)
		if err != nil {
			panic(err)
		}
		set = append(set, uint32(v))
	}
	return set
}

func int32Array(args []string) (set []int32) {
	for _, i := range args {
		v, err := strconv.ParseInt(i, 0, 32)
		if err != nil {
			panic(err)
		}
		set = append(set, int32(v))
	}
	return set
}

func uint64Array(args []string) (set []uint64) {
	for _, u := range args {
		v, err := strconv.ParseUint(u, 0, 64)
		if err != nil {
			panic(err)
		}
		set = append(set, v)
	}
	return set
}

func uint256Array(args []string) (set []*big.Int) {
	for _, u := range args {
		set = append(set, w3.I(u))
	}
	return set
}

func int256Array(args []string) (set []*big.Int) {
	for _, i := range args {
		set = append(set, w3.I(i))
	}
	return set
}

func ParseEncoding(function *w3.Func, signature string, params []string) ([]byte, error) {
	signatureArgs := strings.Split(strings.Split(strings.Split(signature, "(")[1], ")")[0], ",")
	if len(signatureArgs) != len(params) {
		return nil, libcommon.StringError(errors.New("executor parseParams: mismatched arguments"))
	}
	args := []interface{}{}
	for i, s := range signatureArgs {
		var subArgs []string
		if strings.HasSuffix(s, "[]") {
			// set args to an array of s split by commas excluding brackets
			subArgs = strings.Split(strings.Trim(params[i], "[]"), ",")
		}
		switch s {
		case "address":
			args = append(args, w3.A(params[i]))
		case "address[]":
			args = append(args, addressArray(subArgs))
		case "bool":
			args = append(args, params[i] == "true" || params[i] == "TRUE")
		case "bool[]":
			args = append(args, boolArray(subArgs))
		case "string":
			args = append(args, params[i])
		case "string[]":
			args = append(args, stringArray(subArgs))
		case "bytes":
			args = append(args, w3.B(params[i]))
		case "bytes[]":
			args = append(args, bytesArray(subArgs))
		case "uint8":
			v, err := strconv.ParseUint(params[i], 0, 8)
			if err != nil {
				return nil, libcommon.StringError(err)
			}
			args = append(args, v)
		case "uint8[]":
			args = append(args, uint8Array(subArgs))
		case "uint32":
			v, err := strconv.ParseUint(params[i], 0, 32)
			if err != nil {
				return nil, libcommon.StringError(err)
			}
			args = append(args, v)
		case "uint32[]":
			args = append(args, uint32Array(subArgs))
		case "uint256":
			args = append(args, w3.I(params[i]))
		case "uint256[]":
			args = append(args, uint256Array(subArgs))
		case "int8":
			v, err := strconv.ParseInt(params[i], 0, 8)
			if err != nil {
				return nil, libcommon.StringError(err)
			}
			args = append(args, v)
		case "int8[]":
			args = append(args, int8Array(subArgs))
		case "int32":
			v, err := strconv.ParseInt(params[i], 0, 32)
			if err != nil {
				return nil, libcommon.StringError(err)
			}
			args = append(args, v)
		case "int32[]":
			args = append(args, int32Array(subArgs))
		case "int256":
			args = append(args, w3.I(params[i]))
		case "int256[]":
			args = append(args, int256Array(subArgs))
		default:
			return nil, libcommon.StringError(errors.New("executor: parseParams: unsupported type"))
		}
	}
	result, err := function.EncodeArgs(args...)
	if err != nil {
		return nil, libcommon.StringError(err)
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

	if !ValidAddress(addr) {
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

func IsContract(addr string) bool {
	return !IsWallet(addr)
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

func ValidAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}
