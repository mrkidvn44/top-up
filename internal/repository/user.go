package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByPhoneNumber(phoneNumber string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
