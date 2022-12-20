package service

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/store"
)

func GetStringIds(repos repository.Repositories, redis store.RedisStore) (InternalIds, error) {
	empty := InternalIds{}
	ids, err := store.GetObjectFromCache[InternalIds](redis, "internal_ids")
	if err != nil {
		return empty, common.StringError(err)
	}
	if ids == empty {
		user, err := repos.User.GetByType("Internal")
		if err != nil {
			return empty, common.StringError(err)
		}

		bank, err := repos.Instrument.GetBankByUserId(user.ID)
		if err != nil {
			return empty, common.StringError(err)
		}

		wallet, err := repos.Instrument.GetWalletByUserId(user.ID)
		if err != nil {
			return empty, common.StringError(err)
		}

		platform, err := repos.Platform.GetByApiKey("Internal")
		if err != nil {
			return empty, common.StringError(err)
		}

		ids = InternalIds{
			StringUserId:     user.ID,
			StringPlatformId: platform.ID, // temporary
			StringBankId:     bank.ID,
			StringWalletId:   wallet.ID,
		}

		err = store.PutObjectInCache(redis, "internal_ids", ids)

		if err != nil {
			return empty, common.StringError(err)
		}
	}
	return ids, nil
}
