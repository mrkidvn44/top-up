package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type ISkuRepository interface {
	GetSkusByProviderCode(ctx context.Context, providerCode string) (*[]model.Sku, error)
	GetSkuByID(ctx context.Context, id uint) (*model.Sku, error)
	GetSkus(ctx context.Context) (*[]model.Sku, error)
}

type SkuRepository struct {
	db *gorm.DB
}

var _ (ISkuRepository) = (*SkuRepository)(nil)

func NewSkuRepository(db *gorm.DB) *SkuRepository {
	return &SkuRepository{db: db}
}

func (r *SkuRepository) GetSkusByProviderCode(ctx context.Context, providerCode string) (*[]model.Sku, error) {
	var skus []model.Sku
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("Provider").Where("provider_code = ?", providerCode).Find(&skus).Error; err != nil {
		return nil, err
	}
	return &skus, nil
}

func (r *SkuRepository) GetSkuByID(ctx context.Context, id uint) (*model.Sku, error) {
	var sku model.Sku
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("Provider").First(&sku, id).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *SkuRepository) GetSkus(ctx context.Context) (*[]model.Sku, error) {
	var skus []model.Sku
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("Provider").Find(&skus).Error; err != nil {
		return nil, err
	}
	return &skus, nil
}
