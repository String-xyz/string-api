package common

import (
	"crypto/ecdsa"
	"errors"
	"os"
	"strconv"

	"github.com/String-xyz/go-lib/common"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func EVMSign(buffer []byte, eip131 bool) (string, error) {
	privateKey, err := DecryptBlobFromKMS(os.Getenv("EVM_PRIVATE_KEY"))
	if err != nil {
		return "", common.StringError(err)
	}
	return EVMSignWithPrivateKey(buffer, privateKey, eip131)
}

func EVMSignWithPrivateKey(buffer []byte, privateKey string, eip131 bool) (string, error) {
	sk, err := crypto.ToECDSA(ethcommon.FromHex(privateKey))
	if err != nil {
		return "", common.StringError(err)
	}

	if eip131 {
		prefix := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(buffer)))
		buffer = append(prefix, buffer...)
	}

	hash := crypto.Keccak256Hash(buffer)
	signature, err := crypto.Sign(hash.Bytes(), sk)
	if err != nil {
		return "", common.StringError(err)
	}
	return hexutil.Encode(signature), nil
}

func ValidateEVMSignature(signature string, buffer []byte, eip131 bool) (bool, error) {
	// Get private key
	skStr, err := DecryptBlobFromKMS(os.Getenv("EVM_PRIVATE_KEY"))
	if err != nil {
		return false, common.StringError(err)
	}
	sk, err := crypto.ToECDSA(ethcommon.FromHex(skStr))
	if err != nil {
		return false, common.StringError(err)
	}
	pk := sk.Public()
	pkECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return false, common.StringError(errors.New("ValidateSignature: Failed to cast pk to ECDSA"))
	}
	pkBytes := crypto.FromECDSAPub(pkECDSA)

	if eip131 {
		// prepend expected prefix
		prefix := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(buffer)))
		buffer = append(prefix, buffer...)
	}

	hash := crypto.Keccak256Hash(buffer)

	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, common.StringError(err)
	}

	// Handle cases where EIP-155 is not implemented, as with most wallets
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	verified := crypto.VerifySignature(pkBytes, hash.Bytes(), sigBytes[:len(sigBytes)-1]) // last byte of signature is recovery ID
	return verified, nil
}

func ValidateExternalEVMSignature(signature string, address string, buffer []byte, eip131 bool) (bool, error) {
	if eip131 {
		// prepend expected prefix
		prefix := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(buffer)))
		buffer = append(prefix, buffer...)
	}

	hash := crypto.Keccak256Hash(buffer)

	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, common.StringError(err)
	}

	// Handle cases where EIP-155 is not implemented, as with most wallets
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	sigPKECDSA, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false, common.StringError(err)
	}
	sigPKBytes := crypto.FromECDSAPub(sigPKECDSA)

	SIGPKString := crypto.PubkeyToAddress(*sigPKECDSA).String()
	if address != SIGPKString {
		return false, nil
	}

	verified := crypto.VerifySignature(sigPKBytes, hash.Bytes(), sigBytes[:len(sigBytes)-1]) // last byte of signature is recovery ID
	return verified, nil
}
