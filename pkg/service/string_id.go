package service

import (
	"os"
)

func GetStringIdsFromEnv() InternalIds {
	return InternalIds{
		StringUserId:     os.Getenv("STRING_INTERNAL_ID"),
		StringBankId:     os.Getenv("STRING_BANK_ID"),
		StringWalletId:   os.Getenv("STRING_WALLET_ID"),
		StringPlatformId: os.Getenv("STRING_PLACEHOLDER_PLATFORM_ID"), // This is a temporary placeholder, we will get this from API key
	}
}
