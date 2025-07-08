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
	service service.ISkuService
	logger  logger.Interface
}

func NewSkuRouter(handler *gin.RouterGroup, s service.ISkuService, l logger.Interface) {
	h := &SkuRouter{service: s, logger: l}
	skuRoutes := handler.Group("/sku")
	{
		skuRoutes.GET("/:providerCode", h.GetSkusByProviderCode)
		skuRoutes.GET("/", h.GetSkusGroupByProvider)
	}
}

// BasePath /v1/api

// @Summary Get card details by provider code
// @Description Get card details by provider code
// @Tags card-detail
// @Param providerCode path string true "Provider code"
// @Success 200 {array} top-up-api_internal_schema.SkuResponse
// @Router /card-detail/{providerCode} [get]
func (h *SkuRouter) GetSkusByProviderCode(c *gin.Context) {
	providerCode := c.Param("providerCode")
	skus, err := h.service.GetSkusByProviderCode(c, providerCode)
	if err != nil {
		h.logger.Error(errors.New("error getting card details"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(skus))
}

// @Summary Get card details grouped by provider
// @Description Get card details grouped by provider
// @Tags card-detail
// @Success 200 {object} top-up-api_internal_schema.SkusGroupByProvider
// @Router /card-detail [get]
func (h *SkuRouter) GetSkusGroupByProvider(c *gin.Context) {
	skus, err := h.service.GetSkusGroupByProvider(c)
	if err != nil {
		h.logger.Error(errors.New("error getting card details grouped by provider"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if skus == nil {
		c.JSON(http.StatusOK, mapper.SuccessResponse(nil))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(skus))
}
