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

type CardDetailRouter struct {
	service service.CardDetailService
	logger  logger.Interface
}

func NewCardDetailRouter(handler *gin.RouterGroup, s service.CardDetailService, l logger.Interface) {
	h := &CardDetailRouter{service: s, logger: l}
	cardDetailRoutes := handler.Group("/card-detail")
	{
		cardDetailRoutes.GET("/:providerCode", h.GetCardDetailsByProviderCode)
	}
}

// BasePath /v1/api

// @Summary Get card details by provider code
// @Description Get card details by provider code
// @Tags card-detail
// @Param providerCode path string true "Provider code"
// @Success 200 {array} top-up-api_internal_schema.CardDetailResponse
// @Router /card-detail/{providerCode} [get]
func (h *CardDetailRouter) GetCardDetailsByProviderCode(c *gin.Context) {
	providerCode := c.Param("providerCode")
	cardDetails, err := h.service.GetCardDetailsByProviderCode(c, providerCode)
	if err != nil {
		h.logger.Error(errors.New("error getting card details"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schema.SuccessResponse(cardDetails))
}
