package service

import "github.com/String-xyz/string-api/env"

func GetStringIdsFromEnv() InternalIds {
	return InternalIds{
		StringUserId:   env.Var.STRING_INTERNAL_ID,
		StringBankId:   env.Var.STRING_BANK_ID,
		StringWalletId: env.Var.STRING_WALLET_ID,
	}
}
