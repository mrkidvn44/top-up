package service

import (
	"context"
	"errors"
	"testing"
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	"top-up-api/pkg/auth"
	mockRepo "top-up-api/tests/repository/mock"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserService_Login(t *testing.T) {
	mockRepo := new(mockRepo.UserRepositoryMock)
	service := service.NewUserService(mockRepo)
	ctx := context.Background()

	hashed, _ := auth.HashPassword("correct_password")

	tests := []struct {
		name          string
		inputPhone    string
		inputPassword string
		mockUser      *model.User
		mockError     error
		expectedError string
	}{
		{
			name:          "Successful login",
			inputPhone:    "0123456789",
			inputPassword: "correct_password",
			mockUser: &model.User{
				Model:       gorm.Model{ID: 1},
				PhoneNumber: "0123456789",
				Password:    hashed,
			},
			expectedError: "",
		},
		{
			name:          "Wrong password",
			inputPhone:    "0123456789",
			inputPassword: "wrong_password",
			mockUser: &model.User{
				Model:       gorm.Model{ID: 1},
				PhoneNumber: "0123456789",
				Password:    hashed,
			},
			expectedError: "wrong password",
		},
		{
			name:          "User not found",
			inputPhone:    "0000000000",
			inputPassword: "any",
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.On("GetUserByPhoneNumber", tt.inputPhone).Return(tt.mockUser, tt.mockError)

			result, err := service.Login(ctx, tt.inputPhone, tt.inputPassword)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser.ID, result.ID)
				assert.Equal(t, tt.mockUser.PhoneNumber, result.PhoneNumber)
			}
		})
	}
}
func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := new(mockRepo.UserRepositoryMock)
	service := service.NewUserService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name          string
		inputID       uint
		mockUser      *model.User
		mockError     error
		expectedError string
	}{
		{
			name:    "User found",
			inputID: 1,
			mockUser: &model.User{
				Model:       gorm.Model{ID: 1},
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "0123456789",
				Balance:     100,
				CashBack:    10,
			},
			expectedError: "",
		},
		{
			name:          "User not found",
			inputID:       2,
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.On("GetUserByID", ctx, tt.inputID).Return(tt.mockUser, tt.mockError)

			result, err := service.GetUserByID(ctx, tt.inputID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser.ID, result.ID)
				assert.Equal(t, tt.mockUser.FirstName, result.FirstName)
				assert.Equal(t, tt.mockUser.LastName, result.LastName)
				assert.Equal(t, tt.mockUser.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, tt.mockUser.Balance, result.Balance)
				assert.Equal(t, tt.mockUser.CashBack, result.CashBack)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mockRepo.UserRepositoryMock)
	service := service.NewUserService(mockRepo)
	ctx := context.Background()

	validReq := schema.UserCreateRequest{
		FirstName:   "Jane",
		LastName:    "Smith",
		PhoneNumber: "0987654321",
		Password:    "password123",
	}
	hashed, _ := auth.HashPassword(validReq.Password)
	modelUser := &model.User{
		FirstName:   validReq.FirstName,
		LastName:    validReq.LastName,
		PhoneNumber: validReq.PhoneNumber,
		Password:    hashed,
	}

	tests := []struct {
		name      string
		req       schema.UserCreateRequest
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Success",
			req:  validReq,
			mockSetup: func() {
				mockRepo.ExpectedCalls = nil
				mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(u *model.User) bool {
					return u.FirstName == modelUser.FirstName &&
						u.LastName == modelUser.LastName &&
						u.PhoneNumber == modelUser.PhoneNumber &&
						u.Password != "1"
				})).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "User already exists (pg error)",
			req:  validReq,
			mockSetup: func() {
				mockRepo.ExpectedCalls = nil
				mockRepo.On("CreateUser", ctx, mock.Anything).Return(&pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				})
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.CreateUser(ctx, tt.req)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
