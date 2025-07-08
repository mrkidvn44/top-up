package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type ProviderService interface {
	GetProviders(ctx context.Context) (*[]schema.ProviderResponse, error)
}

type providerService struct {
	repo repository.ProviderRepository
}

var _ ProviderService = (*providerService)(nil)

func NewProviderService(providerRepository repository.ProviderRepository) *providerService {
	return &providerService{repo: providerRepository}
}

func (s *providerService) GetProviders(ctx context.Context) (*[]schema.ProviderResponse, error) {
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
