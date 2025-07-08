package mock

import (
	"context"
	"top-up-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type PurchaseHistoryRepositoryMock struct {
	mock.Mock
}

func (m *PurchaseHistoryRepositoryMock) CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error {
	args := m.Called(ctx, purchaseHistory)
	return args.Error(0)
}

func (m *PurchaseHistoryRepositoryMock) GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) ([]model.PurchaseHistory, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Error(2) != nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}

	return args.Get(0).([]model.PurchaseHistory), args.Get(1).(int64), args.Error(2)
}

func (m *PurchaseHistoryRepositoryMock) GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.PurchaseHistory), args.Error(1)
}

func (m *PurchaseHistoryRepositoryMock) UpdatePurchaseHistoryStatusByOrderID(ctx context.Context, order_id uint, status model.PurchaseHistoryStatus) error {
	args := m.Called(ctx, order_id, status)
	return args.Error(0)
}

func (m *PurchaseHistoryRepositoryMock) GetPurchaseHistoryByOrderID(ctx context.Context, order_id uint) (*model.PurchaseHistory, error) {
	args := m.Called(ctx, order_id)
	return args.Get(0).(*model.PurchaseHistory), args.Error(1)
}
