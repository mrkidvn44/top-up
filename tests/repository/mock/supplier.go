package mock

import (
	"context"
	"top-up-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type SupplierRepositoryMock struct {
	mock.Mock
}

func (m *SupplierRepositoryMock) GetSuppliers(ctx context.Context) (*[]model.Supplier, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Supplier), args.Error(1)
}
