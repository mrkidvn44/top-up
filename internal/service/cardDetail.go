package service

import (
	"context"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type CardDetailService struct {
	repo repository.CardDetailRepository
}

func NewCardDetailService(repo repository.CardDetailRepository) *CardDetailService {
	return &CardDetailService{repo: repo}
}

func (s *CardDetailService) GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]schema.CardDetailResponse, error) {
	return s.repo.GetCardDetailsByProviderCode(ctx, providerCode)
}
