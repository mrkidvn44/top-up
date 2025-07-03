package redis

import (
	"context"
	"time"
	"top-up-api/config"

	"github.com/redis/go-redis/v9"
)

const (
	_timeOut = 5 * time.Minute
)

var NotFound = redis.Nil

type Interface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	GetLock(ctx context.Context, key string) bool
	ReleaseLock(ctx context.Context, key string) error
	TryAcquireLock(ctx context.Context, key string, timeout time.Duration) error
}

type RedisClient struct {
	Client *redis.Client
}

var _ Interface = (*RedisClient)(nil)

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

func (r *RedisClient) GetLock(ctx context.Context, key string) bool {
	encodeKey := "lock:" + key
	wasSet, err := r.Client.SetNX(ctx, encodeKey, 1, _timeOut).Result()
	return err == nil && wasSet
}

func (r *RedisClient) ReleaseLock(ctx context.Context, key string) error {
	encodeKey := "lock:" + key
	return r.Client.Del(ctx, encodeKey).Err()
}

// TryAcquireLock tries to acquire a lock for the given key within the specified timeout duration.
// Returns nil if lock acquired, otherwise returns error if timed out.
func (r *RedisClient) TryAcquireLock(ctx context.Context, key string, timeout time.Duration) error {
	expireTime := time.Now().Add(timeout)
	for {
		if ok := r.GetLock(ctx, key); ok {
			return nil
		}
		if time.Now().After(expireTime) {
			return context.DeadlineExceeded
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// ...existing code...
