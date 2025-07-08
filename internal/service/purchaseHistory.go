package service

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/model"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
)

type PurchaseHistoryService interface {
	GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) (*schema.PaginationResponse, error)
	GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error)
}

type purchaseHistoryService struct {
	repo repository.PurchaseHistoryRepository
}

var _ PurchaseHistoryService = (*purchaseHistoryService)(nil)

func NewPurchaseHistoryService(repo repository.PurchaseHistoryRepository) *purchaseHistoryService {
	return &purchaseHistoryService{repo: repo}
}

func (s *purchaseHistoryService) GetPurchaseHistoriesByUserIDPaginated(ctx context.Context, userID uint, page, pageSize int) (*schema.PaginationResponse, error) {
	histories, total, err := s.repo.GetPurchaseHistoriesByUserIDPaginated(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	historiesResponse := make([]*schema.PurchaseHistoryResponse, len(histories))
	for i, history := range histories {
		historiesResponse[i] = mapper.PurchaseHistoryResponseFromModel(&history)
	}

	totalPage := (int(total) + pageSize - 1) / pageSize
	return mapper.PaginationResponseFromModel(
		int(total),
		totalPage,
		page,
		historiesResponse,
	), nil
}

func (s *purchaseHistoryService) GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error) {
	return s.repo.GetPurchaseHistoryByID(ctx, id)
}
