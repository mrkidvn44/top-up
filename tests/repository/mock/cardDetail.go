package mock

import (
	"context"
	"top-up-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type CardDetailRepositoryMock struct {
	mock.Mock
}

func (m *CardDetailRepositoryMock) GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]model.CardDetail, error) {
	args := m.Called(ctx, providerCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.CardDetail), args.Error(1)
}

func (m *CardDetailRepositoryMock) GetCardDetailByID(ctx context.Context, id uint) (*model.CardDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CardDetail), args.Error(1)
}

func (m *CardDetailRepositoryMock) GetCardDetails(ctx context.Context) (*[]model.CardDetail, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.CardDetail), args.Error(1)
}
