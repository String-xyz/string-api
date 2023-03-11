package store

import (
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
)

func NewRedis() database.RedisStore {
	opts := database.RedisConfigOptions{
		Host:        os.Getenv("REDIS_HOST"),
		Port:        os.Getenv("REDIS_PORT"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		ClusterMode: !common.IsLocalEnv(),
	}
	return database.NewRedisStore(opts)
}
