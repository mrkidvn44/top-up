package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type SkuService interface {
	GetSkusBySupplierCode(ctx context.Context, supplierCode string) (*[]schema.SkuResponse, error)
	GetSkusGroupBySupplier(ctx context.Context) (*[]schema.SkusGroupBySupplier, error)
}

type skuService struct {
	repo repository.SkuRepository
}

var _ SkuService = (*skuService)(nil)

func NewSkuService(repo repository.SkuRepository) *skuService {
	return &skuService{repo: repo}
}

func (s *skuService) GetSkusBySupplierCode(ctx context.Context, supplierCode string) (*[]schema.SkuResponse, error) {
	skus, err := s.repo.GetSkusBySupplierCode(ctx, supplierCode)
	if err != nil {
		return nil, err
	}
	if skus == nil {
		return nil, nil
	}
	skuResponses := make([]schema.SkuResponse, len(*skus))
	for i, sku := range *skus {
		skuResponses[i] = *mapper.SkuResponseFromModel(sku)
	}
	return &skuResponses, nil
}

func (s *skuService) GetSkusGroupBySupplier(ctx context.Context) (*[]schema.SkusGroupBySupplier, error) {
	skus, err := s.repo.GetSkus(ctx)
	if err != nil {
		return nil, err
	}
	if skus == nil {
		return nil, nil
	}
	groupedDetails := mapper.SkusGroupBySupplierFromModel(*skus)
	if groupedDetails == nil {
		return nil, nil
	}

	return groupedDetails, nil
}
