package controller

import (
	"errors"
	"payment-api/internal/service"
	"payment-api/pkg/logger"
	orderpb "payment-api/proto/order"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type OrderRouter struct {
	service    service.IOrderService
	logger     logger.Interface
	grpcClient orderpb.OrderServiceClient
}

func NewOrderRouter(handler *gin.RouterGroup, s service.IOrderService, l logger.Interface, grpcClient orderpb.OrderServiceClient) {
	h := &OrderRouter{service: s, logger: l, grpcClient: grpcClient}
	orderRoutes := handler.Group("/order")
	{
		orderRoutes.POST("/create", h.CreateOrder)
		orderRoutes.PATCH("/update", h.UpdateOrder)
	}
}

func (h *OrderRouter) CreateOrder(c *gin.Context) {

	data, err := c.GetRawData()
	if err != nil {
		h.logger.Error(errors.New("failed to read request body"), zap.Error(err))
		return
	}

	h.service.CreateOrder(string(data), h.grpcClient)

	c.JSON(200, gin.H{"message": "Order created successfully"})

}

func (h *OrderRouter) UpdateOrder(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		h.logger.Error(errors.New("failed to read request body"), zap.Error(err))
		return
	}
	h.logger.Info("Received order update data", zapcore.Field{
		Key:    "data",
		Type:   zapcore.StringType,
		String: string(data),
	})
	c.JSON(200, gin.H{"message": "Order updated successfully"})
}
