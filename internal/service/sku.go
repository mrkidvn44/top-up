package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type SkuService interface {
	GetSkusByProviderCode(ctx context.Context, providerCode string) (*[]schema.SkuResponse, error)
	GetSkusGroupByProvider(ctx context.Context) (*[]schema.SkusGroupByProvider, error)
}

type skuService struct {
	repo repository.SkuRepository
}

var _ SkuService = (*skuService)(nil)

func NewSkuService(repo repository.SkuRepository) *skuService {
	return &skuService{repo: repo}
}

func (s *skuService) GetSkusByProviderCode(ctx context.Context, providerCode string) (*[]schema.SkuResponse, error) {
	skus, err := s.repo.GetSkusByProviderCode(ctx, providerCode)
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

func (s *skuService) GetSkusGroupByProvider(ctx context.Context) (*[]schema.SkusGroupByProvider, error) {
	skus, err := s.repo.GetSkus(ctx)
	if err != nil {
		return nil, err
	}
	if skus == nil {
		return nil, nil
	}
	groupedDetails := mapper.SkusGroupByProviderFromModel(*skus)
	if groupedDetails == nil {
		return nil, nil
	}

	return groupedDetails, nil
}
