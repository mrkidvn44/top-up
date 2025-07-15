package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func (m *RedisMock) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *RedisMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *RedisMock) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *RedisMock) ReleaseLock(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *RedisMock) TryAcquireLock(ctx context.Context, key string, timeout time.Duration) error {
	args := m.Called(ctx, key, timeout)
	return args.Error(0)
}
