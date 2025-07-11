package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type ProviderRepository interface {
	GetProvidersWithSuppliers(ctx context.Context) ([]model.Provider, error)
}

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) GetProvidersWithSuppliers(ctx context.Context) ([]model.Provider, error) {
	var providers []model.Provider
	err := r.db.WithContext(ctx).Preload("Suppliers").Find(&providers).Error
	return providers, err
}
