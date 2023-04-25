package service

import "github.com/String-xyz/string-api/env"

func GetStringIdsFromEnv() InternalIds {
	user, _ := env.Get("STRING_INTERNAL_ID")
	bank, _ := env.Get("STRING_BANK_ID")
	wallet, _ := env.Get("STRING_WALLET_ID")
	return InternalIds{
		StringUserId:   user,
		StringBankId:   bank,
		StringWalletId: wallet,
	}
}
