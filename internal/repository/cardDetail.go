package repository

import (
	"context"
	"top-up-api/internal/model"
	"top-up-api/internal/schema"

	"gorm.io/gorm"
)

type CardDetailRepository struct {
	db *gorm.DB
}

func NewCardDetailRepository(db *gorm.DB) *CardDetailRepository {
	return &CardDetailRepository{db: db}
}

func (r *CardDetailRepository) GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]schema.CardDetailResponse, error) {
	var cardDetails []model.CardDetail
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("CardPrice").Preload("Provider").Where("provider_code = ?", providerCode).Find(&cardDetails).Error; err != nil {
		return nil, err
	}
	cardDetailResponses := make([]schema.CardDetailResponse, len(cardDetails))
	for i, cardDetail := range cardDetails {
		cardDetailResponses[i] = *schema.CardDetailResponseFromModel(cardDetail)
	}
	return &cardDetailResponses, nil
}

func (r *CardDetailRepository) GetCardDetailByID(ctx context.Context, id uint) (*model.CardDetail, error) {
	var cardDetail model.CardDetail
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("CardPrice").Preload("Provider").First(&cardDetail, id).Error; err != nil {
		return nil, err
	}
	return &cardDetail, nil
}
