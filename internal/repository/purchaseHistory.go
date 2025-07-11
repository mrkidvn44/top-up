package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type PurchaseHistoryRepository interface {
	CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error
	GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) ([]model.PurchaseHistory, int64, error)
	GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error)
	UpdatePurchaseHistoryStatusByOrderID(ctx context.Context, order_id uint, status model.PurchaseHistoryStatus) error
	GetPurchaseHistoryByOrderID(ctx context.Context, order_id uint) (*model.PurchaseHistory, error)
}

type purchaseHistoryRepository struct {
	db *gorm.DB
}

var _ PurchaseHistoryRepository = (*purchaseHistoryRepository)(nil)

func NewPurchaseHistoryRepository(db *gorm.DB) *purchaseHistoryRepository {
	return &purchaseHistoryRepository{db: db}
}

func (r *purchaseHistoryRepository) CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error {
	return r.db.WithContext(ctx).Create(purchaseHistory).Error
}

func (r *purchaseHistoryRepository) GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) ([]model.PurchaseHistory, int64, error) {
	var histories []model.PurchaseHistory
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model.PurchaseHistory{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Sku").
		Preload("Sku.Supplier").
		Preload("Sku.CashBack").
		Limit(pageSize).
		Offset(offset).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *purchaseHistoryRepository) GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error) {
	var purchaseHistory model.PurchaseHistory
	if err := r.db.WithContext(ctx).Where("order_id = ?", id).First(&purchaseHistory).Error; err != nil {
		return nil, err
	}
	return &purchaseHistory, nil
}

func (r *purchaseHistoryRepository) UpdatePurchaseHistoryStatusByOrderID(ctx context.Context, order_id uint, status model.PurchaseHistoryStatus) error {
	return r.db.WithContext(ctx).Model(&model.PurchaseHistory{}).
		Where("order_id = ?", order_id).
		Update("status", status).Error
}

func (r *purchaseHistoryRepository) GetPurchaseHistoryByOrderID(ctx context.Context, order_id uint) (*model.PurchaseHistory, error) {
	var purchaseHistory model.PurchaseHistory
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", order_id).
		Preload("User").
		Preload("Sku").
		Preload("Sku.Supplier").
		Preload("Sku.CashBack").
		First(&purchaseHistory).Error; err != nil {
		return nil, err
	}
	return &purchaseHistory, nil
}
