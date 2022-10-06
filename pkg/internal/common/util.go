package common

import (
	"crypto/sha256"
	"encoding/hex"

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
		return ethcomm.Address{}, err
	}
	return crypto.PubkeyToAddress(*recovered), nil
}
