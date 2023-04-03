package service

import (
	"os"
)

func GetStringIdsFromEnv() InternalIds {
	return InternalIds{
		StringUserId:   os.Getenv("STRING_INTERNAL_ID"),
		StringBankId:   os.Getenv("STRING_BANK_ID"),
		StringWalletId: os.Getenv("STRING_WALLET_ID"),
	}
}
