package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type IPurchaseHistoryRepository interface {
	CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error
	GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) ([]model.PurchaseHistory, int64, error)
	GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error)
	UpdatePurchaseHistoryStatusByOrderID(ctx context.Context, order_id uint, status model.PurchaseHistoryStatus) error
	GetPurchaseHistoryByOrderID(ctx context.Context, order_id uint) (*model.PurchaseHistory, error)
}

type PurchaseHistoryRepository struct {
	db *gorm.DB
}

var _ IPurchaseHistoryRepository = (*PurchaseHistoryRepository)(nil)

func NewPurchaseHistoryRepository(db *gorm.DB) *PurchaseHistoryRepository {
	return &PurchaseHistoryRepository{db: db}
}

func (r *PurchaseHistoryRepository) CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error {
	return r.db.WithContext(ctx).Create(purchaseHistory).Error
}

func (r *PurchaseHistoryRepository) GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) ([]model.PurchaseHistory, int64, error) {
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
		Preload("CardDetail").
		Preload("CardDetail.Provider").
		Preload("CardDetail.CardPrice").
		Preload("CardDetail.CashBack").
		Limit(pageSize).
		Offset(offset).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *PurchaseHistoryRepository) GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error) {
	var purchaseHistory model.PurchaseHistory
	if err := r.db.WithContext(ctx).Where("order_id = ?", id).First(&purchaseHistory).Error; err != nil {
		return nil, err
	}
	return &purchaseHistory, nil
}

func (r *PurchaseHistoryRepository) UpdatePurchaseHistoryStatusByOrderID(ctx context.Context, order_id uint, status model.PurchaseHistoryStatus) error {
	return r.db.WithContext(ctx).Model(&model.PurchaseHistory{}).
		Where("order_id = ?", order_id).
		Update("status", status).Error
}

func (r *PurchaseHistoryRepository) GetPurchaseHistoryByOrderID(ctx context.Context, order_id uint) (*model.PurchaseHistory, error) {
	var purchaseHistory model.PurchaseHistory
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", order_id).
		Preload("User").
		Preload("CardDetail").
		Preload("CardDetail.Provider").
		Preload("CardDetail.CardPrice").
		Preload("CardDetail.CashBack").
		First(&purchaseHistory).Error; err != nil {
		return nil, err
	}
	return &purchaseHistory, nil
}
