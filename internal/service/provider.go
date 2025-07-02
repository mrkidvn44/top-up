package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type IProviderService interface {
	GetProviders(ctx context.Context) (*[]schema.ProviderResponse, error)
}

type ProviderService struct {
	repo repository.IProviderRepository
}

func NewProviderService(providerRepository repository.IProviderRepository) *ProviderService {
	return &ProviderService{repo: providerRepository}
}

func (s *ProviderService) GetProviders(ctx context.Context) (*[]schema.ProviderResponse, error) {
	providers, err := s.repo.GetProviders(ctx)
	if err != nil {
		return nil, err
	}
	providerResponses := make([]schema.ProviderResponse, len(*providers))
	for i, provider := range *providers {
		providerResponses[i] = *mapper.ProviderResponseFromModel(&provider)
	}
	return &providerResponses, nil
}
