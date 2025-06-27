package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

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

func (r *UserRepository) DeductBalance(ctx context.Context, userID uint, amount int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance - ?", amount)).Error
}

func (r *UserRepository) AddBalance(ctx context.Context, userID uint, amount int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *UserRepository) AddCashBack(ctx context.Context, userID uint, amount int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("cash_back", gorm.Expr("cash_back + ?", amount)).Error
}
