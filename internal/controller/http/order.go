package controller

import (
	"errors"
	"net/http"
	"strconv"
	"top-up-api/internal/mapper"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/logger"
	"top-up-api/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderRouter struct {
	service   service.IOrderService
	logger    logger.Interface
	auth      auth.Interface
	validator validator.Interface
}

func NewOrderRouter(handler *gin.RouterGroup, s service.IOrderService, l logger.Interface, a auth.Interface, v validator.Interface) {
	h := &OrderRouter{service: s, logger: l, auth: a, validator: v}
	orderRoutes := handler.Group("/order")
	{
		orderRoutes.POST("/create", h.CreateOrder)
		orderRoutes.POST("/confirm", h.ConfirmOrder)
		orderRoutes.PATCH("/update-status", h.UpdateOrderStatus)
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
// @Security Bearer
func (h *OrderRouter) CreateOrder(c *gin.Context) {
	orderRequest := schema.OrderRequest{}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		h.logger.Error(errors.New("failed to bind order request"), zap.Error(err))
		c.JSON(http.StatusBadRequest, mapper.ErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	token, err := h.auth.AuthenticateService(c)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusUnauthorized, mapper.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	_, err = h.auth.ValidateUserIDFromToken(token, strconv.Itoa(int(orderRequest.UserID)))
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusUnauthorized, mapper.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	orderResponse, err := h.service.CreateOrder(c, orderRequest)
	if err != nil {
		h.logger.Error(errors.New("failed to create order"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, mapper.ErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(orderResponse))
}

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
		c.JSON(http.StatusBadRequest, mapper.ErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	
	if err := h.validator.Validate(orderConfirmRequest); err != nil {
		h.logger.Error(errors.New("validation failed for order confirm request"), zap.Error(err))
		c.JSON(http.StatusBadRequest, mapper.ErrorResponse(http.StatusBadRequest, "Validation Error", err.Error()))
		return
	}

	err := h.service.ConfirmOrder(c, orderConfirmRequest)
	if err != nil {
		h.logger.Error(errors.New("failed to confirm order"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, mapper.ErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(nil))
}

// Summary Update order status
// @Description Update order status
// @Tags order
// @Accept json
// @Produce json
// @Param orderUpdateRequest body top-up-api_internal_schema.OrderUpdateRequest true "Order update request"
// @Success 200 {object} top-up-api_internal_schema.Response
// @Router /order/update-status [patch]
func (h *OrderRouter) UpdateOrderStatus(c *gin.Context) {
	orderUpdateRequest := schema.OrderUpdateRequest{}
	if err := c.ShouldBindJSON(&orderUpdateRequest); err != nil {
		h.logger.Error(errors.New("failed to bind order update request"), zap.Error(err))
		c.JSON(http.StatusBadRequest, mapper.ErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	err := h.service.UpdateOrderStatus(c, orderUpdateRequest)
	if err != nil {
		h.logger.Error(errors.New("failed to update order status"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, mapper.ErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(nil))
}
