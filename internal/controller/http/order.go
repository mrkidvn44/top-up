package controller

import (
	"errors"
	"net/http"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	"top-up-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderRouter struct {
	service service.OrderService
	logger  logger.Interface
}

func NewOrderRouter(handler *gin.RouterGroup, s service.OrderService, l logger.Interface) {
	h := &OrderRouter{service: s, logger: l}
	orderRoutes := handler.Group("/order")
	{
		orderRoutes.POST("/create", h.CreateOrder)
		orderRoutes.POST("/confirm", h.ConfirmOrder)
	}
}

// BasePath /v1/api

// @Summary Create order
// @Description Create order
// @Tags order
// @Accept json
// @Produce json
// @Param orderRequest body top-up-api_internal_schema.OrderRequest true "Order request"
// @Success 200 {object} top-up-api_internal_schema.OrderResponse
// @Router /order/create [post]
func (h *OrderRouter) CreateOrder(c *gin.Context) {
	orderRequest := schema.OrderRequest{}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		h.logger.Error(errors.New("failed to bind order request"), zap.Error(err))
		c.JSON(http.StatusBadRequest, schema.ErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	orderResponse, err := h.service.CreateOrder(c, orderRequest)
	if err != nil {
		h.logger.Error(errors.New("failed to create order"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, schema.SuccessResponse(orderResponse))
}

// BasePath /v1/api

// @Summary Confirm order
// @Description Confirm order
// @Tags order
// @Accept json
// @Produce json
// @Param orderConfirmRequest body top-up-api_internal_schema.OrderConfirmRequest true "Order confirm request"
// @Success 200 {object} top-up-api_internal_schema.OrderConfirmRequest
// @Router /order/confirm [post]
func (h *OrderRouter) ConfirmOrder(c *gin.Context) {
	orderConfirmRequest := schema.OrderConfirmRequest{}
	if err := c.ShouldBindJSON(&orderConfirmRequest); err != nil {
		h.logger.Error(errors.New("failed to bind order confirm request"), zap.Error(err))
		c.JSON(http.StatusBadRequest, schema.ErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	err := h.service.ConfirmOrder(c, orderConfirmRequest)
	if err != nil {
		h.logger.Error(errors.New("failed to confirm order"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, schema.SuccessResponse(nil))
}
