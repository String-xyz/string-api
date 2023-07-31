package service

import "github.com/String-xyz/string-api/config"

func GetStringIdsFromEnv() InternalIds {
	return InternalIds{
		StringUserId:   config.Var.STRING_INTERNAL_ID,
		StringBankId:   config.Var.STRING_BANK_ID,
		StringWalletId: config.Var.STRING_WALLET_ID,
	}
}
