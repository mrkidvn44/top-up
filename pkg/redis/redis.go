package redis

import (
	"context"
	"time"
	"top-up-api/config"

	"github.com/redis/go-redis/v9"
)

var NotFound = redis.Nil

type RedisClient struct {
	Client *redis.Client
}

type Interface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

func NewRedis(cfg config.Redis) *RedisClient {
	return &RedisClient{Client: redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
