package service

import (
	"context"
	"top-up-api/internal/model"
	"top-up-api/internal/repository"
)

type PurchaseHistoryService struct {
	repo repository.PurchaseHistoryRepository
}

func NewPurchaseHistoryService(repo repository.PurchaseHistoryRepository) *PurchaseHistoryService {
	return &PurchaseHistoryService{repo: repo}
}

func (s *PurchaseHistoryService) GetPurchaseHistoryByUserID(ctx context.Context, userID uint) (*model.PurchaseHistory, error) {
	return s.repo.GetPurchaseHistoryByUserID(ctx, userID)
}

func (s *PurchaseHistoryService) CreatePurchaseHistory(ctx context.Context, purchaseHistory *model.PurchaseHistory) error {
	return s.repo.CreatePurchaseHistory(ctx, purchaseHistory)
}

func (s *PurchaseHistoryService) GetPurchaseHistoryByID(ctx context.Context, id uint) (*model.PurchaseHistory, error) {
	return s.repo.GetPurchaseHistoryByID(ctx, id)
}

func (s *PurchaseHistoryService) UpdatePurchaseHistoryStatus(ctx context.Context, id uint, status model.PurchaseHistoryStatus) error {
	return s.repo.UpdatePurchaseHistoryStatus(ctx, id, status)
}
