package rediscache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
)

type RedisCache struct {
	c *redis.Client
}

func NewRedisCache(
	e *config.Env,
) *RedisCache {
	opts, err := redis.ParseURL(e.RedisDatabaseURL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}

	return &RedisCache{
		c: client,
	}
}

func (r *RedisCache) Scan(
	ctx context.Context,
	key cache.Key,
	value any,
) (bool, error) {
	strCmd := r.c.Get(ctx, string(key))
	if strCmd.Err() == redis.Nil {
		return false, nil
	}
	if strCmd.Err() != nil {
		return false, strCmd.Err()
	}

	if err := strCmd.Scan(value); err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisCache) Set(
	ctx context.Context,
	key cache.Key,
	value any,
	expiration time.Duration,
) error {
	return r.c.Set(ctx, string(key), value, expiration).Err()
}

func (r *RedisCache) Delete(
	ctx context.Context,
	keys ...cache.Key,
) error {
	ks := make([]string, len(keys))
	for i, k := range keys {
		ks[i] = string(k)
	}

	return r.c.Del(ctx, ks...).Err()
}

var _ cache.Cache = &RedisCache{}
