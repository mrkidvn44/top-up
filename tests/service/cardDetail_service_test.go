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

func TestCardDetailService_GetCardDetailsByProviderCode(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name              string
		providerCode      string
		setupMock         func(*mockRepo.CardDetailRepositoryMock)
		expectedDetails   *[]model.CardDetail
		expectedError     string
		expectedNilResult bool
	}{
		{
			name:         "Success",
			providerCode: "VTL",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				cardDetails := &[]model.CardDetail{
					{ProviderCode: "VTL", CardPriceCode: "CP1"},
					{ProviderCode: "VTL", CardPriceCode: "CP2"},
				}
				m.ExpectedCalls = nil
				m.On("GetCardDetailsByProviderCode", ctx, "VTL").Return(cardDetails, nil)
			},
			expectedDetails: &[]model.CardDetail{
				{ProviderCode: "VTL", CardPriceCode: "CP1"},
				{ProviderCode: "VTL", CardPriceCode: "CP2"},
			},
			expectedError:     "",
			expectedNilResult: false,
		},
		{
			name:         "Repository error",
			providerCode: "VTL",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetCardDetailsByProviderCode", ctx, "VTL").Return(nil, errors.New("db error"))
			},
			expectedDetails:   nil,
			expectedError:     "db error",
			expectedNilResult: true,
		},
		{
			name:         "No record",
			providerCode: "VTL",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				empty := &[]model.CardDetail{}
				m.ExpectedCalls = nil
				m.On("GetCardDetailsByProviderCode", ctx, "VTL").Return(empty, nil)
			},
			expectedDetails:   &[]model.CardDetail{},
			expectedError:     "",
			expectedNilResult: false,
		},
		{
			name:         "Nil from repo",
			providerCode: "VTL",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetCardDetailsByProviderCode", ctx, "VTL").Return(nil, nil)
			},
			expectedDetails:   nil,
			expectedError:     "",
			expectedNilResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.CardDetailRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewCardDetailService(mockRepo)
			got, err := svc.GetCardDetailsByProviderCode(ctx, tt.providerCode)
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

func TestCardDetailService_GetCardDetailsGroupByProvider(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		setupMock      func(*mockRepo.CardDetailRepositoryMock)
		expectedResult *[]model.CardDetail
		expectedError  string
		expectedNil    bool
	}{
		{
			name: "Success",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				cardDetails := &[]model.CardDetail{
					{
						ProviderCode:  "VTL",
						CardPriceCode: "CP1",
						Provider:      model.Provider{Code: "VTL", Name: "Viettel", LogoUrl: "logo1"},
					},
					{
						ProviderCode:  "MBF",
						CardPriceCode: "CP2",
						Provider:      model.Provider{Code: "MBF", Name: "Mobifone", LogoUrl: "logo2"},
					},
				}
				m.ExpectedCalls = nil
				m.On("GetCardDetails", ctx).Return(cardDetails, nil)
			},
			expectedResult: nil,
			expectedError:  "",
			expectedNil:    false,
		},
		{
			name: "Repository error",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetCardDetails", ctx).Return(nil, errors.New("db error"))
			},
			expectedResult: nil,
			expectedError:  "db error",
			expectedNil:    true,
		},
		{
			name: "No record",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				empty := &[]model.CardDetail{}
				m.ExpectedCalls = nil
				m.On("GetCardDetails", ctx).Return(empty, nil)
			},
			expectedResult: nil,
			expectedError:  "",
			expectedNil:    true,
		},
		{
			name: "Nil from repo",
			setupMock: func(m *mockRepo.CardDetailRepositoryMock) {
				m.ExpectedCalls = nil
				m.On("GetCardDetails", ctx).Return(nil, nil)
			},
			expectedResult: nil,
			expectedError:  "",
			expectedNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo.CardDetailRepositoryMock)
			tt.setupMock(mockRepo)
			svc := service.NewCardDetailService(mockRepo)
			got, err := svc.GetCardDetailsGroupByProvider(ctx)
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
					if group.ProviderCode == "VTL" {
						foundVTL = true
					}
					if group.ProviderCode == "MBF" {
						foundMBF = true
					}
				}
				assert.True(t, foundVTL)
				assert.True(t, foundMBF)
			}
		})
	}
}
