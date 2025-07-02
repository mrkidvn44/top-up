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

func TestProviderService_GetProviders(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name              string
		setupMock         func(*mockRepo.ProviderRepositoryMock)
		args              args
		expectedProviders *[]model.Provider
		expectedError     string
	}{
		{
			name: "Success",
			setupMock: func(m *mockRepo.ProviderRepositoryMock) {
				providers := &[]model.Provider{
					{
						Code:    "VTL",
						Name:    "Viettel",
						LogoUrl: "viettel.png",
						Status:  string(model.ProviderStatusActive),
					},
					{
						Code:    "MBF",
						Name:    "Mobifone",
						LogoUrl: "mobifone.png",
						Status:  string(model.ProviderStatusActive),
					},
				}
				m.ExpectedCalls = nil
				m.On("GetProviders", ctx).Return(providers, nil)
			},
			args: args{ctx: ctx},
			expectedProviders: &[]model.Provider{
				{
					Code:    "VTL",
					Name:    "Viettel",
					LogoUrl: "viettel.png",
					Status:  string(model.ProviderStatusActive),
				},
				{
					Code:    "MBF",
					Name:    "Mobifone",
					LogoUrl: "mobifone.png",
					Status:  string(model.ProviderStatusActive),
				},
			},
			expectedError: "",
		},
		{
			name: "Repository error",
			setupMock: func(m *mockRepo.ProviderRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetProviders", ctx).Return(nil, errors.New("db error"))
			},
			args:              args{ctx: ctx},
			expectedProviders: nil,
			expectedError:     "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.ProviderRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewProviderService(mockRepo)
			got, err := svc.GetProviders(tt.args.ctx)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, len(*tt.expectedProviders), len(*got))
				for i := range *tt.expectedProviders {
					assert.Equal(t, (*tt.expectedProviders)[i].Code, (*got)[i].Code)
					assert.Equal(t, (*tt.expectedProviders)[i].Name, (*got)[i].Name)
				}
			}
		})
	}
}
