package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	mockRepo "top-up-api/tests/repository/mock"
	"top-up-api/tests/util"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Predefined mock data and expected responses
var (
	ctx = context.Background()

	mockPurchaseHistory1 = model.PurchaseHistory{
		Model:         gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		OrderID:       1001,
		UserID:        1,
		SkuID:         2001,
		TotalPrice:    10000,
		PhoneNumber:   "081234567890",
		Status:        model.PurchaseHistoryStatusSuccess,
		CashBackValue: 100,
		Sku:           *util.CreateMockSku(2001, "VTL", 10000, model.CashBackTypeFixed, 0, "Viettel"),
	}
	mockPurchaseHistory2 = model.PurchaseHistory{
		Model:         gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		OrderID:       1002,
		UserID:        1,
		SkuID:         2002,
		TotalPrice:    20000,
		PhoneNumber:   "081234567891",
		Status:        model.PurchaseHistoryStatusConfirm,
		CashBackValue: 200,
		Sku:           *util.CreateMockSku(2002, "MBF", 20000, model.CashBackTypeFixed, 0, "Mobifone"),
	}

	mockPaginationEmpty = &schema.PaginationResponse{
		Code:       200,
		Message:    "success",
		Data:       []*schema.PurchaseHistoryResponse{},
		Pagination: schema.Pagination{TotalCount: 0, TotalPage: 0, CurrentPage: 1},
	}
	mockPaginationWithData = &schema.PaginationResponse{
		Code:       200,
		Message:    "success",
		Data:       []*schema.PurchaseHistoryResponse{}, // Data checked in test
		Pagination: schema.Pagination{TotalCount: 2, TotalPage: 1, CurrentPage: 1},
	}
	mockPaginationCalc = &schema.PaginationResponse{
		Code:       200,
		Message:    "success",
		Data:       []*schema.PurchaseHistoryResponse{},
		Pagination: schema.Pagination{TotalCount: 12, TotalPage: 3, CurrentPage: 2},
	}

	expectedPurchaseHistory = &model.PurchaseHistory{
		Model:         gorm.Model{ID: 1},
		OrderID:       1001,
		UserID:        1,
		SkuID:         2001,
		TotalPrice:    10000,
		PhoneNumber:   "081234567890",
		Status:        model.PurchaseHistoryStatusSuccess,
		CashBackValue: 100,
	}
)

func TestPurchaseHistoryService_GetPurchaseHistoriesByUserIDPaginated(t *testing.T) {
	type args struct {
		ctx      context.Context
		userID   uint
		page     int
		pageSize int
	}

	tests := []struct {
		name          string
		setupMock     func(*mockRepo.PurchaseHistoryRepositoryMock)
		args          args
		expectedResp  *schema.PaginationResponse
		expectedError string
	}{
		{
			name: "Success - empty data",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoriesByUserIDPaginated", ctx, uint(1), 1, 10).Return([]model.PurchaseHistory{}, int64(0), nil)
			},
			args:          args{ctx: ctx, userID: 1, page: 1, pageSize: 10},
			expectedResp:  mockPaginationEmpty,
			expectedError: "",
		},
		{
			name: "Success - with data",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoriesByUserIDPaginated", ctx, uint(1), 1, 10).Return([]model.PurchaseHistory{mockPurchaseHistory1, mockPurchaseHistory2}, int64(2), nil)
			},
			args:          args{ctx: ctx, userID: 1, page: 1, pageSize: 10},
			expectedResp:  mockPaginationWithData,
			expectedError: "",
		},
		{
			name: "Success - pagination calculation test",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoriesByUserIDPaginated", ctx, uint(1), 2, 5).Return([]model.PurchaseHistory{}, int64(12), nil)
			},
			args:          args{ctx: ctx, userID: 1, page: 2, pageSize: 5},
			expectedResp:  mockPaginationCalc,
			expectedError: "",
		},
		{
			name: "Repository error",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoriesByUserIDPaginated", ctx, uint(1), 1, 10).Return(nil, int64(0), errors.New("db error"))
			},
			args:          args{ctx: ctx, userID: 1, page: 1, pageSize: 10},
			expectedResp:  nil,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewPurchaseHistoryService(mockRepo)
			got, err := svc.GetPurchaseHistoriesByUserIDPaginated(tt.args.ctx, tt.args.userID, tt.args.page, tt.args.pageSize)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.expectedResp.Pagination.TotalPage, got.Pagination.TotalPage)
				assert.Equal(t, tt.expectedResp.Pagination.TotalCount, got.Pagination.TotalCount)
				assert.Equal(t, tt.expectedResp.Pagination.CurrentPage, got.Pagination.CurrentPage)

				actualData, ok := got.Data.([]*schema.PurchaseHistoryResponse)
				assert.True(t, ok, "Data should be of type []*schema.PurchaseHistoryResponse")

				if tt.name == "Success - with data" {
					assert.Len(t, actualData, 2, "Should have 2 purchase history records")
					assert.Equal(t, uint(1001), actualData[0].OrderID)
					assert.Equal(t, uint(1), actualData[0].UserID)
					assert.Equal(t, uint(2001), actualData[0].SkuID)
					assert.Equal(t, 10000, actualData[0].TotalPrice)
					assert.Equal(t, "081234567890", actualData[0].PhoneNumber)
					assert.Equal(t, "success", actualData[0].Status)
					assert.Equal(t, 100, actualData[0].CashBackValue)
				} else {
					assert.Len(t, actualData, 0, "Should have empty data")
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPurchaseHistoryService_GetPurchaseHistoryByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint
	}

	tests := []struct {
		name          string
		setupMock     func(*mockRepo.PurchaseHistoryRepositoryMock)
		args          args
		expectedResp  *model.PurchaseHistory
		expectedError string
	}{
		{
			name: "Success",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoryByID", ctx, uint(1)).Return(&mockPurchaseHistory1, nil)
			},
			args:          args{ctx: ctx, id: 1},
			expectedResp:  expectedPurchaseHistory,
			expectedError: "",
		},
		{
			name: "Repository error",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoryByID", ctx, uint(999)).Return((*model.PurchaseHistory)(nil), errors.New("record not found"))
			},
			args:          args{ctx: ctx, id: 999},
			expectedResp:  nil,
			expectedError: "record not found",
		},
		{
			name: "Not found",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoryByID", ctx, uint(404)).Return((*model.PurchaseHistory)(nil), gorm.ErrRecordNotFound)
			},
			args:          args{ctx: ctx, id: 404},
			expectedResp:  nil,
			expectedError: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewPurchaseHistoryService(mockRepo)
			got, err := svc.GetPurchaseHistoryByID(tt.args.ctx, tt.args.id)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.expectedResp.ID, got.ID)
				assert.Equal(t, tt.expectedResp.OrderID, got.OrderID)
				assert.Equal(t, tt.expectedResp.UserID, got.UserID)
				assert.Equal(t, tt.expectedResp.SkuID, got.SkuID)
				assert.Equal(t, tt.expectedResp.TotalPrice, got.TotalPrice)
				assert.Equal(t, tt.expectedResp.PhoneNumber, got.PhoneNumber)
				assert.Equal(t, tt.expectedResp.Status, got.Status)
				assert.Equal(t, tt.expectedResp.CashBackValue, got.CashBackValue)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
