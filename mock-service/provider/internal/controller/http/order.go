package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"provider-api/internal/service"
	"provider-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderRouter struct {
	service service.IOrderService
	logger  logger.Interface
}

func NewOrderRouter(handler *gin.RouterGroup, s service.IOrderService, l logger.Interface) {
	h := &OrderRouter{service: s, logger: l}
	orderRoutes := handler.Group("/order")
	{
		orderRoutes.POST("/", h.ProcessOrder)
	}
}

func (h *OrderRouter) ProcessOrder(c *gin.Context) {
	var orderData struct {
		OrderID     uint   `json:"order_id"`
		PhoneNumber string `json:"phone_number"`
		TotalPrice  int    `json:"total_price"`
		CardPrice   int    `json:"card_price"`
		CallBackURL string `json:"callback_url"`
	}

	if err := c.ShouldBindJSON(&orderData); err != nil {
		h.logger.Error(errors.New("failed to bind JSON"), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	fmt.Print("Received order data: ", orderData)
	updateRequest := struct {
		OrderID uint   `json:"order_id"`
		Status  string `json:"status"`
	}{
		OrderID: orderData.OrderID,
		Status:  h.service.ProcessOrder(),
	}

	go func(updateRequest interface{}) {
		payload, _ := json.Marshal(updateRequest)
		req, err := http.NewRequest(http.MethodPatch, orderData.CallBackURL, bytes.NewBuffer(payload))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
	}(updateRequest)

	c.JSON(200, gin.H{"message": "Order update sent successfully", "order_id": orderData.OrderID, "status": updateRequest.Status})
}
