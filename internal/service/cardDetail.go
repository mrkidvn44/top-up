package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type ISkuService interface {
	GetSkusByProviderCode(ctx context.Context, providerCode string) (*[]schema.SkuResponse, error)
	GetSkusGroupByProvider(ctx context.Context) (*[]schema.SkusGroupByProvider, error)
}

type SkuService struct {
	repo repository.ISkuRepository
}

var _ ISkuService = (*SkuService)(nil)

func NewSkuService(repo repository.ISkuRepository) *SkuService {
	return &SkuService{repo: repo}
}

func (s *SkuService) GetSkusByProviderCode(ctx context.Context, providerCode string) (*[]schema.SkuResponse, error) {
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

func (s *SkuService) GetSkusGroupByProvider(ctx context.Context) (*[]schema.SkusGroupByProvider, error) {
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
