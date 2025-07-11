package controller

import (
	"errors"
	"net/http"
	"top-up-api/internal/mapper"
	"top-up-api/internal/service"
	"top-up-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SkuRouter struct {
	service service.SkuService
	logger  logger.Interface
}

func NewSkuRouter(handler *gin.RouterGroup, s service.SkuService, l logger.Interface) {
	h := &SkuRouter{service: s, logger: l}
	skuRoutes := handler.Group("/sku")
	{
		skuRoutes.GET("/:supplierCode", h.GetSkusBySupplierCode)
		skuRoutes.GET("/", h.GetSkusGroupBySupplier)
	}
}

// BasePath /v1/api

// @Summary Get sku details by supplier code
// @Description Get sku details by supplier code
// @Tags sku
// @Param supplierCode path string true "Supplier code"
// @Success 200 {array} top-up-api_internal_schema.SkuResponse
// @Router /sku/{supplierCode} [get]
func (h *SkuRouter) GetSkusBySupplierCode(c *gin.Context) {
	supplierCode := c.Param("supplierCode")
	skus, err := h.service.GetSkusBySupplierCode(c, supplierCode)
	if err != nil {
		h.logger.Error(errors.New("error getting card details"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(skus))
}

// @Summary Get card details grouped by supplier
// @Description Get card details grouped by supplier
// @Tags sku
// @Success 200 {object} top-up-api_internal_schema.SkusGroupBySupplier
// @Router /sku [get]
func (h *SkuRouter) GetSkusGroupBySupplier(c *gin.Context) {
	skus, err := h.service.GetSkusGroupBySupplier(c)
	if err != nil {
		h.logger.Error(errors.New("error getting card details grouped by supplier"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if skus == nil {
		c.JSON(http.StatusOK, mapper.SuccessResponse(nil))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(skus))
}
