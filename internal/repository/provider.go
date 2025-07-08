package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type ProviderRepository interface {
	GetProviders(ctx context.Context) (*[]model.Provider, error)
}
type providerRepository struct {
	db *gorm.DB
}

var _ ProviderRepository = (*providerRepository)(nil)

func NewProviderRepository(db *gorm.DB) *providerRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) GetProviders(ctx context.Context) (*[]model.Provider, error) {
	var providers []model.Provider
	if err := r.db.WithContext(ctx).Find(&providers).Error; err != nil {
		return nil, err
	}
	return &providers, nil
}
