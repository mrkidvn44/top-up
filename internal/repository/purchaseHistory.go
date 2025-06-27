package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type PurchaseHistoryRepository struct {
	db *gorm.DB
}

func NewPurchaseHistoryRepository(db *gorm.DB) *PurchaseHistoryRepository {
	return &PurchaseHistoryRepository{db: db}
}

func (r *PurchaseHistoryRepository) CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error {
	return r.db.WithContext(ctx).Create(purchaseHistory).Error
}

func (r *PurchaseHistoryRepository) GetPurchaseHistoryByUserID(ctx context.Context, userID uint) (*model.PurchaseHistory, error) {
	var purchaseHistory model.PurchaseHistory
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&purchaseHistory).Error; err != nil {
		return nil, err
	}
	return &purchaseHistory, nil
}

func (r *PurchaseHistoryRepository) GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error) {
	var purchaseHistory model.PurchaseHistory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&purchaseHistory).Error; err != nil {
		return nil, err
	}
	return &purchaseHistory, nil
}

func (r *PurchaseHistoryRepository) UpdatePurchaseHistoryStatus(ctx context.Context, id uint, status model.PurchaseHistoryStatus) error {
	return r.db.WithContext(ctx).Model(&model.PurchaseHistory{}).
		Where("id = ?", id).
		Update("status", status).Error
}
