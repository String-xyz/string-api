package service

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateWallet() error {
	overwriteMe, _ := GetAddress()
	if overwriteMe != "" {
		fmt.Printf("\n WARNING!  Generating a new wallet will overwrite the old one.  To prevent loss of funds, the old wallet credentials will be printed below:")
		sk, _ := GetPrivateKey()
		fmt.Printf("\n ADDRESS: %+v\nPRIVATE KEY: %+v", overwriteMe, sk)
	}

	sk, err := crypto.GenerateKey()
	if err != nil {
		return common.StringError(err)
	}

	blob, err := common.EncryptStringToKMS(hexutil.Encode(crypto.FromECDSA(sk))[2:])
	if err != nil {
		return common.StringError(err)
	}
	err = common.PutSSM("string-encrypted-sk", string(blob), true)
	if err != nil {
		return common.StringError(err)
	}

	pk := sk.Public()
	pkECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return common.StringError(err)
	}
	addrStr := crypto.PubkeyToAddress(*pkECDSA).Hex()
	fmt.Printf("\nGenerated and encrypted new private key to SSM for wallet %+v", addrStr)

	err = common.PutSSM("string-wallet-address", addrStr, true)
	if err != nil {
		return common.StringError(err)
	}
	fmt.Printf("\nWallet address was also added to SSM.")
	return nil
}

func GetPrivateKey() (string, error) {
	blobStr, err := common.GetSSM("string-encrypted-sk")
	if err != nil {
		return "", common.StringError(err)
	}
	sk, err := common.DecryptBlobFromKMS(blobStr)
	if err != nil {
		return "", common.StringError(err)
	}
	return sk, nil
}

func GetAddress() (string, error) {
	address, err := common.GetSSM("string-wallet-address")
	if err != nil {
		return "", common.StringError(err)
	}
	return address, nil
}
