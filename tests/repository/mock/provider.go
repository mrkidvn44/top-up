package mock

import (
	"context"
	"top-up-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type ProviderRepositoryMock struct {
	mock.Mock
}

func (m *ProviderRepositoryMock) GetProvidersWithSuppliers(ctx context.Context) ([]model.Provider, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Provider), args.Error(1)
}
