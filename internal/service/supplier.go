package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type SupplierService interface {
	GetSuppliers(ctx context.Context) (*[]schema.SupplierResponse, error)
}

type supplierService struct {
	repo repository.SupplierRepository
}

var _ SupplierService = (*supplierService)(nil)

func NewSupplierService(supplierRepository repository.SupplierRepository) *supplierService {
	return &supplierService{repo: supplierRepository}
}

func (s *supplierService) GetSuppliers(ctx context.Context) (*[]schema.SupplierResponse, error) {
	suppliers, err := s.repo.GetSuppliers(ctx)
	if err != nil {
		return nil, err
	}
	supplierResponses := make([]schema.SupplierResponse, len(*suppliers))
	for i, supplier := range *suppliers {
		supplierResponses[i] = *mapper.SupplierResponseFromModel(&supplier)
	}
	return &supplierResponses, nil
}
