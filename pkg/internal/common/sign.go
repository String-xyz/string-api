package common

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func EVMSign(data interface{}) (string, error) {
	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	if err != nil {
		return "", StringError(err)
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		return "", StringError(err)
	}
	hash := crypto.Keccak256Hash(buffer)
	signature, err := crypto.Sign(hash.Bytes(), sk)
	if err != nil {
		return "", StringError(err)
	}
	fmt.Printf("\nSIGNED=%+v", hexutil.Encode(signature))
	return hexutil.Encode(signature), nil
}

func ValidateEVMSignature(signature string, data interface{}) (bool, error) {
	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("EVM_PRIVATE_KEY")))
	if err != nil {
		return false, StringError(err)
	}
	pk := sk.Public()
	pkECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return false, StringError(errors.New("ValidateSignature: Failed to cast pk to ECDSA"))
	}
	pkBytes := crypto.FromECDSAPub(pkECDSA)

	buffer, err := json.Marshal(data)
	if err != nil {
		return false, StringError(err)
	}
	hash := crypto.Keccak256Hash(buffer)

	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, StringError(err)
	}
	verified := crypto.VerifySignature(pkBytes, hash.Bytes(), sigBytes[:len(sigBytes)-1]) // last byte of signature is recovery ID
	return verified, nil
}

func ValidateExternalEVMSignature(signature string, address string, data interface{}) (bool, error) {
	buffer, err := json.Marshal(data)
	if err != nil {
		return false, StringError(err)
	}
	hash := crypto.Keccak256Hash(buffer)

	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, StringError(err)
	}

	sigPKECDSA, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false, StringError(err)
	}
	sigPKBytes := crypto.FromECDSAPub(sigPKECDSA)

	SIGPKString := crypto.PubkeyToAddress(*sigPKECDSA).String()
	if address != SIGPKString {
		return false, nil
	}

	verified := crypto.VerifySignature(sigPKBytes, hash.Bytes(), sigBytes[:len(sigBytes)-1]) // last byte of signature is recovery ID
	return verified, nil
}
