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
	ReleaseLock(ctx context.Context, key string) error
	TryAcquireLock(ctx context.Context, key string, timeout time.Duration) error
}

type redisClient struct {
	Client *redis.Client
}

var _ Interface = (*redisClient)(nil)

func NewRedis(cfg config.Redis) *redisClient {
	return &redisClient{Client: redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})}
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *redisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *redisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *redisClient) getLock(ctx context.Context, encodeKey string) (bool, error) {
	wasSet, err := r.Client.SetNX(ctx, encodeKey, 1, _timeOut).Result()
	return wasSet, err
}

func (r *redisClient) ReleaseLock(ctx context.Context, key string) error {
	encodeKey := getEncodeKey(key)
	releaseChannel := getReleashKey(encodeKey)
	err := r.Client.Del(ctx, encodeKey).Err()
	if err != nil {
		return err
	}

	return r.Client.Publish(ctx, releaseChannel, "released").Err()
}

func (r *redisClient) TryAcquireLock(ctx context.Context, key string, timeout time.Duration) error {
	expireTime := time.Now().Add(timeout)
	encodeKey := getEncodeKey(key)
	releaseChannel := getReleashKey(encodeKey)

	pubsub := r.Client.Subscribe(ctx, releaseChannel)
	defer pubsub.Close()

	for {
		ok, err := r.getLock(ctx, encodeKey)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Until(expireTime)):
			return context.DeadlineExceeded
		case <-pubsub.Channel():
		}
	}
}

func getEncodeKey(key string) string {
	return "lock:" + key
}

func getReleashKey(encodeKey string) string {
	return encodeKey + ":release"
}
