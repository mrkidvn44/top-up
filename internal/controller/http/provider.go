package controller

import (
	"errors"
	"net/http"
	"top-up-api/internal/mapper"
	service "top-up-api/internal/service"
	"top-up-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProviderRouter struct {
	service service.IProviderService
	logger  logger.Interface
}

func NewProviderRouter(handler *gin.RouterGroup, s service.IProviderService, l logger.Interface) {
	h := &ProviderRouter{service: s, logger: l}
	providerRoutes := handler.Group("/provider")
	{
		providerRoutes.GET("/", h.GetProviders)
	}
}

// BasePath /v1/api

// @Summary Get providers
// @Description Get providers
// @Tags provider
// @Accept json
// @Produce json
// @Success 200 {array} top-up-api_internal_schema.ProviderResponse
// @Router /provider [get]
func (h *ProviderRouter) GetProviders(c *gin.Context) {
	providers, err := h.service.GetProviders(c)
	if err != nil {
		h.logger.Error(errors.New("error getting providers"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, mapper.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(providers))
}
