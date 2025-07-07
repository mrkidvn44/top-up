package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type ICardDetailService interface {
	GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]schema.CardDetailResponse, error)
	GetCardDetailsGroupByProvider(ctx context.Context) (*[]schema.CardDetailsGroupByProvider, error)
}

type CardDetailService struct {
	repo repository.ICardDetailRepository
}

var _ ICardDetailService = (*CardDetailService)(nil)

func NewCardDetailService(repo repository.ICardDetailRepository) *CardDetailService {
	return &CardDetailService{repo: repo}
}

func (s *CardDetailService) GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]schema.CardDetailResponse, error) {
	cardDetails, err := s.repo.GetCardDetailsByProviderCode(ctx, providerCode)
	if err != nil {
		return nil, err
	}
	if cardDetails == nil {
		return nil, nil
	}
	cardDetailResponses := make([]schema.CardDetailResponse, len(*cardDetails))
	for i, cardDetail := range *cardDetails {
		cardDetailResponses[i] = *mapper.CardDetailResponseFromModel(cardDetail)
	}
	return &cardDetailResponses, nil
}

func (s *CardDetailService) GetCardDetailsGroupByProvider(ctx context.Context) (*[]schema.CardDetailsGroupByProvider, error) {
	cardDetails, err := s.repo.GetCardDetails(ctx)
	if err != nil {
		return nil, err
	}
	if cardDetails == nil {
		return nil, nil
	}
	groupedDetails := mapper.CardDetailsGroupByProviderFromModel(*cardDetails)
	if groupedDetails == nil {
		return nil, nil
	}

	return groupedDetails, nil
}
