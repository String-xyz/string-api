package store

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisStore interface {
	Get(id string) ([]byte, error)
	Set(string, any, time.Duration) error
	HSet(string, map[string]interface{}) error
	HGetAll(string) (map[string]string, error)
	HDel(string, string) int64
	HMLen(string) int64
	Delete(string) error
}

type redisStore struct {
	client *redis.Client
}

func redisConf() *redis.Options {
	url := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	var tlsCf *tls.Config
	if os.Getenv("ENV") != "local" {
		tlsCf = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	cf := &redis.Options{
		Addr:      url,
		TLSConfig: tlsCf,
		Password:  os.Getenv("REDIS_PASSWORD"),
		DB:        0,
	}
	return cf
}

func NewRedisStore() RedisStore {
	ctx := context.Background()
	client := redis.NewClient(redisConf())
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	return &redisStore{
		client: client,
	}
}

func (r redisStore) Delete(id string) error {
	_, err := r.client.Del(r.client.Context(), id).Result()
	if err != nil {
		return common.StringError(errors.Wrap(err, "unable to delete"))
	}
	return nil
}

func (r redisStore) Get(id string) ([]byte, error) {
	ctx := context.Background()
	return r.client.Get(ctx, id).Bytes()
}

func (r redisStore) Set(id string, value any, expire time.Duration) error {
	ctx := context.Background()
	if err := r.client.Set(ctx, id, value, expire).Err(); err != nil {
		return common.StringError(errors.Wrap(err, "failed to save value to redis"))
	}
	return nil
}

func (r redisStore) HSet(key string, data map[string]interface{}) error {
	ctx := context.Background()
	if err := r.client.HSet(ctx, key, data).Err(); err != nil {
		return common.StringError(errors.Wrap(err, "failed to save array to redis"))
	}

	return nil
}

func (r redisStore) HGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return data, common.StringError(err)
	}
	return data, nil
}

func (r redisStore) HMLen(key string) int64 {
	ctx := context.Background()
	data := r.client.HLen(ctx, key)
	return data.Val()
}

func (r redisStore) HDel(key, val string) int64 {
	ctx := context.Background()
	data := r.client.HDel(ctx, key, val)
	return data.Val()
}
