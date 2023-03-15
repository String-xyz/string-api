package store

import (
	"os"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
)

func NewRedis() database.RedisStore {
	opts := database.RedisConfigOptions{
		Host:        os.Getenv("REDIS_HOST"),
		Port:        os.Getenv("REDIS_PORT"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		ClusterMode: !libcommon.IsLocalEnv(),
	}
	return database.NewRedisStore(opts)
}
