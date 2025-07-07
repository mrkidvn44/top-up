package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	mockRepo "top-up-api/tests/repository/mock"

	"github.com/stretchr/testify/assert"
)

func TestPurchaseHistoryService_GetPurchaseHistoriesByUserIDPaginated(t *testing.T) {
	ctx := context.Background()
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
			name: "Success",
			setupMock: func(m *mockRepo.PurchaseHistoryRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetPurchaseHistoriesByUserIDPaginated", ctx, uint(1), 1, 10).Return([]model.PurchaseHistory{}, int64(0), nil)
			},
			args: args{ctx: ctx, userID: 1, page: 1, pageSize: 10},
			expectedResp: &schema.PaginationResponse{
				Code:       200,
				Message:    "success",
				Data:       []model.PurchaseHistory{},
				Pagination: schema.Pagination{TotalCount: 0, TotalPage: 0, CurrentPage: 1},
			},
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
			fmt.Print(err)
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
				expectedData, _ := tt.expectedResp.Data.([]model.PurchaseHistory)
				actualData, _ := got.Data.([]model.PurchaseHistory)
				assert.Equal(t, len(expectedData), len(actualData))
			}
		})
	}
}
