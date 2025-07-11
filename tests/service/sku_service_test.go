package service

import (
	"context"
	"errors"
	"testing"
	"top-up-api/internal/model"
	"top-up-api/internal/service"
	mockRepo "top-up-api/tests/repository/mock"

	"github.com/stretchr/testify/assert"
)

func TestSkuService_GetSkusBySupplierCode(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name              string
		supplierCode      string
		setupMock         func(*mockRepo.SkuRepositoryMock)
		expectedDetails   *[]model.Sku
		expectedError     string
		expectedNilResult bool
	}{
		{
			name:         "Success",
			supplierCode: "VTL",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				skus := &[]model.Sku{
					{SupplierCode: "VTL"},
					{SupplierCode: "VTL"},
				}
				m.ExpectedCalls = nil
				m.On("GetSkusBySupplierCode", ctx, "VTL").Return(skus, nil)
			},
			expectedDetails: &[]model.Sku{
				{SupplierCode: "VTL"},
				{SupplierCode: "VTL"},
			},
			expectedError:     "",
			expectedNilResult: false,
		},
		{
			name:         "Repository error",
			supplierCode: "VTL",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetSkusBySupplierCode", ctx, "VTL").Return(nil, errors.New("db error"))
			},
			expectedDetails:   nil,
			expectedError:     "db error",
			expectedNilResult: true,
		},
		{
			name:         "No record",
			supplierCode: "VTL",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				empty := &[]model.Sku{}
				m.ExpectedCalls = nil
				m.On("GetSkusBySupplierCode", ctx, "VTL").Return(empty, nil)
			},
			expectedDetails:   &[]model.Sku{},
			expectedError:     "",
			expectedNilResult: false,
		},
		{
			name:         "Nil from repo",
			supplierCode: "VTL",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetSkusBySupplierCode", ctx, "VTL").Return(nil, nil)
			},
			expectedDetails:   nil,
			expectedError:     "",
			expectedNilResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.SkuRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewSkuService(mockRepo)
			got, err := svc.GetSkusBySupplierCode(ctx, tt.supplierCode)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, got)
			} else if tt.expectedNilResult {
				assert.NoError(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, len(*tt.expectedDetails), len(*got))
			}
		})
	}
}

func TestSkuService_GetSkusGroupBySupplier(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		setupMock      func(*mockRepo.SkuRepositoryMock)
		expectedResult *[]model.Sku
		expectedError  string
		expectedNil    bool
	}{
		{
			name: "Success",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				skus := &[]model.Sku{
					{
						SupplierCode: "VTL",
						Supplier:     model.Supplier{Code: "VTL", Name: "Viettel", LogoUrl: "logo1"},
					},
					{
						SupplierCode: "MBF",
						Supplier:     model.Supplier{Code: "MBF", Name: "Mobifone", LogoUrl: "logo2"},
					},
				}
				m.ExpectedCalls = nil
				m.On("GetSkus", ctx).Return(skus, nil)
			},
			expectedResult: nil,
			expectedError:  "",
			expectedNil:    false,
		},
		{
			name: "Repository error",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetSkus", ctx).Return(nil, errors.New("db error"))
			},
			expectedResult: nil,
			expectedError:  "db error",
			expectedNil:    true,
		},
		{
			name: "No record",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				empty := &[]model.Sku{}
				m.ExpectedCalls = nil
				m.On("GetSkus", ctx).Return(empty, nil)
			},
			expectedResult: nil,
			expectedError:  "",
			expectedNil:    true,
		},
		{
			name: "Nil from repo",
			setupMock: func(m *mockRepo.SkuRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetSkus", ctx).Return(nil, nil)
			},
			expectedResult: nil,
			expectedError:  "",
			expectedNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.SkuRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewSkuService(mockRepo)
			got, err := svc.GetSkusGroupBySupplier(ctx)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, got)
			} else if tt.expectedNil {
				assert.NoError(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				foundVTL := false
				foundMBF := false
				for _, group := range *got {
					if group.SupplierCode == "VTL" {
						foundVTL = true
					}
					if group.SupplierCode == "MBF" {
						foundMBF = true
					}
				}
				assert.True(t, foundVTL)
				assert.True(t, foundMBF)
			}
		})
	}
}
