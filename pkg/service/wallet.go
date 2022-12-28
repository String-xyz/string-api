package service

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateWallet() error {
	sk, err := crypto.GenerateKey()
	if err != nil {
		return common.StringError(err)
	}

	blob, err := common.EncryptStringToKMS(hexutil.Encode(crypto.FromECDSA(sk))[2:])
	common.PutSSM("STRING_ENCRYPTED_SK", string(blob))

	pk := sk.Public()
	pkECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return common.StringError(err)
	}
	addrStr := crypto.PubkeyToAddress(*pkECDSA).Hex()
	fmt.Printf("\nGenerated and encrypted new private key to SSM: %+v", addrStr)
	return nil
}

func GetPrivateKey() (string, error) {
	blobStr, err := common.GetSSM("STRING_ENCRYPTED_SK")
	if err != nil {
		return "", common.StringError(err)
	}
	sk, err := common.DecryptBlobFromKMS([]byte(blobStr))
	if err != nil {
		return "", common.StringError(err)
	}
	return sk, nil
}
