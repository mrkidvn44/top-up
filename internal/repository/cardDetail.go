package repository

import (
	"context"
	"top-up-api/internal/model"

	"gorm.io/gorm"
)

type ICardDetailRepository interface {
	GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]model.CardDetail, error)
	GetCardDetailByID(ctx context.Context, id uint) (*model.CardDetail, error)
	GetCardDetails(ctx context.Context) (*[]model.CardDetail, error)
}

type CardDetailRepository struct {
	db *gorm.DB
}

var _ (ICardDetailRepository) = (*CardDetailRepository)(nil)

func NewCardDetailRepository(db *gorm.DB) *CardDetailRepository {
	return &CardDetailRepository{db: db}
}

func (r *CardDetailRepository) GetCardDetailsByProviderCode(ctx context.Context, providerCode string) (*[]model.CardDetail, error) {
	var cardDetails []model.CardDetail
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("CardPrice").Preload("Provider").Where("provider_code = ?", providerCode).Find(&cardDetails).Error; err != nil {
		return nil, err
	}
	return &cardDetails, nil
}

func (r *CardDetailRepository) GetCardDetailByID(ctx context.Context, id uint) (*model.CardDetail, error) {
	var cardDetail model.CardDetail
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("CardPrice").Preload("Provider").First(&cardDetail, id).Error; err != nil {
		return nil, err
	}
	return &cardDetail, nil
}

func (r *CardDetailRepository) GetCardDetails(ctx context.Context) (*[]model.CardDetail, error) {
	var cardDetails []model.CardDetail
	if err := r.db.WithContext(ctx).Preload("CashBack").Preload("CardPrice").Preload("Provider").Find(&cardDetails).Error; err != nil {
		return nil, err
	}
	return &cardDetails, nil
}
