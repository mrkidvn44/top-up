package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type IProviderRepository interface {
	GetProviders(ctx context.Context) (*[]model.Provider, error)
}
type ProviderRepository struct {
	db *gorm.DB
}

var _ IProviderRepository = (*ProviderRepository)(nil)

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) GetProviders(ctx context.Context) (*[]model.Provider, error) {
	var providers []model.Provider
	if err := r.db.WithContext(ctx).Find(&providers).Error; err != nil {
		return nil, err
	}
	return &providers, nil
}
