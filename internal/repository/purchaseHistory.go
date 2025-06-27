package repository

import (
	"context"
	"errors"
	"top-up-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (r *PurchaseHistoryRepository) UpdatePurchaseHistoryStatusWithLock(ctx context.Context, orderID uint, status model.PurchaseHistoryStatus) error {
	// Begin a transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Fetch the purchase history with a FOR UPDATE lock
		var purchaseHistory model.PurchaseHistory
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ?", orderID).First(&purchaseHistory).Error; err != nil {
			return errors.New("purchase history not found: " + err.Error())
		}

		// Validate the status
		if purchaseHistory.Status != model.PurchaseHistoryStatusConfirm {
			return errors.New("order is not confirmed")
		}

		// Update the status
		if err := tx.Model(&model.PurchaseHistory{}).Where("order_id = ?", orderID).Update("status", status).Error; err != nil {
			return errors.New("failed to update order status: " + err.Error())
		}

		return nil
	})
}
