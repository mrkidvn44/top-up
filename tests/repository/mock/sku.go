package mock

import (
	"context"
	"top-up-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type SkuRepositoryMock struct {
	mock.Mock
}

func (m *SkuRepositoryMock) GetSkusBySupplierCode(ctx context.Context, supplierCode string) (*[]model.Sku, error) {
	args := m.Called(ctx, supplierCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Sku), args.Error(1)
}

func (m *SkuRepositoryMock) GetSkuByID(ctx context.Context, id uint) (*model.Sku, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Sku), args.Error(1)
}

func (m *SkuRepositoryMock) GetSkus(ctx context.Context) (*[]model.Sku, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Sku), args.Error(1)
}
