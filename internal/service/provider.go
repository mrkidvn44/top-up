package service

import (
	"context"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type ProviderService struct {
	repo repository.ProviderRepository
}

func NewProviderService(providerRepository repository.ProviderRepository) *ProviderService {
	return &ProviderService{repo: providerRepository}
}

func (s *ProviderService) GetProviders(ctx context.Context) (*[]schema.ProviderResponse, error) {
	providers, err := s.repo.GetProviders(ctx)
	if err != nil {
		return nil, err
	}
	providerResponses := make([]schema.ProviderResponse, len(*providers))
	for i, provider := range *providers {
		providerResponses[i] = *schema.ProviderResponseFromModel(&provider)
	}
	return &providerResponses, nil
}
