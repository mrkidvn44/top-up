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

func TestSupplierService_GetSuppliers(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name              string
		setupMock         func(*mockRepo.SupplierRepositoryMock)
		args              args
		expectedSuppliers *[]model.Supplier
		expectedError     string
	}{
		{
			name: "Success",
			setupMock: func(m *mockRepo.SupplierRepositoryMock) {
				Suppliers := &[]model.Supplier{
					{
						Code:    "VTL",
						Name:    "Viettel",
						LogoUrl: "viettel.png",
						Status:  model.SupplierStatusActive,
					},
					{
						Code:    "MBF",
						Name:    "Mobifone",
						LogoUrl: "mobifone.png",
						Status:  model.SupplierStatusActive,
					},
				}
				m.ExpectedCalls = nil
				m.On("GetSuppliers", ctx).Return(Suppliers, nil)
			},
			args: args{ctx: ctx},
			expectedSuppliers: &[]model.Supplier{
				{
					Code:    "VTL",
					Name:    "Viettel",
					LogoUrl: "viettel.png",
					Status:  model.SupplierStatusActive,
				},
				{
					Code:    "MBF",
					Name:    "Mobifone",
					LogoUrl: "mobifone.png",
					Status:  model.SupplierStatusActive,
				},
			},
			expectedError: "",
		},
		{
			name: "Repository error",
			setupMock: func(m *mockRepo.SupplierRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetSuppliers", ctx).Return(nil, errors.New("db error"))
			},
			args:              args{ctx: ctx},
			expectedSuppliers: nil,
			expectedError:     "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.SupplierRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewSupplierService(mockRepo)
			got, err := svc.GetSuppliers(tt.args.ctx)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, len(*tt.expectedSuppliers), len(*got))
				for i := range *tt.expectedSuppliers {
					assert.Equal(t, (*tt.expectedSuppliers)[i].Code, (*got)[i].Code)
					assert.Equal(t, (*tt.expectedSuppliers)[i].Name, (*got)[i].Name)
				}
			}
		})
	}
}
