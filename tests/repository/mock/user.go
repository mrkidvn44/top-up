// internal/repository/mock/user_repository_mock.go
package mock

import (
	"context"
	"top-up-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByPhoneNumber(phone string) (*model.User, error) {
	args := m.Called(phone)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
