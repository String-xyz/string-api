package service

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/String-xyz/string-api/pkg/store"
)

func GetStringUUIDs(repos repository.Repositories, redis store.RedisStore) (InternalUUIDs, error) {
	empty := InternalUUIDs{}
	uuids, err := store.GetObjectFromCache[InternalUUIDs](redis, "internal_uuids")
	if err != nil {
		return empty, common.StringError(err)
	}
	if uuids == empty {
		user, err := repos.User.GetByType("Internal")
		if err != nil {
			return empty, common.StringError(err)
		}

		device, err := repos.Device.GetByUserId(user.ID)
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
		uuids = InternalUUIDs{
			stringUserId:   user.ID,
			stringDeviceId: device.ID,
			// stringPlatformId: "",
			StringBankId:   bank.ID,
			StringWalletId: wallet.ID,
		}
		store.PutObjectInCache(redis, "internal_uuids", uuids)
	}
	return uuids, nil
}
