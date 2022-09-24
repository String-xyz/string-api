package common

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func Sign(data interface{}) (string, error) {
	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	if err != nil {
		return "", err
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := crypto.Keccak256Hash(buffer)
	signature, err := crypto.Sign(hash.Bytes(), sk)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(signature), nil
}

func ValidateSignature(signature string, data interface{}) (bool, error) {
	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	if err != nil {
		return false, err
	}
	pk := sk.Public()
	pkECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("ValidateSignature: Failed to cast pk to ECDSA")
	}
	pkBytes := crypto.FromECDSAPub(pkECDSA)

	buffer, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	hash := crypto.Keccak256Hash(buffer)

	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}
	sigBytes = sigBytes[:len(sigBytes)-1] // last byte is a recovery ID
	verified := crypto.VerifySignature(pkBytes, hash.Bytes(), sigBytes)
	return verified, nil
}
