package service

import (
	"math/rand/v2"
	"provider-api/pkg/logger"
	"time"
)

const (
	CacheTime = 30 * time.Minute
	CacheKey  = "order_id"
)

type IOrderService interface {
	ProcessOrder() string
}

type OrderService struct {
	logger logger.Interface
}

var _ IOrderService = (*OrderService)(nil)

func NewOrderService(l logger.Interface) *OrderService {
	return &OrderService{
		logger: l,
	}
}

func (s *OrderService) ProcessOrder() string {
	randInt := rand.IntN(100)
	if randInt < 10 {
		return "failed"
	}
	if randInt < 20 {
		return "pending"
	}
	return "success"
}
