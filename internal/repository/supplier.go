package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type SupplierRepository interface {
	GetSuppliers(ctx context.Context) (*[]model.Supplier, error)
}
type supplierRepository struct {
	db *gorm.DB
}

var _ SupplierRepository = (*supplierRepository)(nil)

func NewSupplierRepository(db *gorm.DB) *supplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) GetSuppliers(ctx context.Context) (*[]model.Supplier, error) {
	var suppliers []model.Supplier
	if err := r.db.WithContext(ctx).Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return &suppliers, nil
}
