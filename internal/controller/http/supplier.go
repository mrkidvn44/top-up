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

type SupplierRouter struct {
	service service.SupplierService
	logger  logger.Interface
}

func NewSupplierRouter(handler *gin.RouterGroup, s service.SupplierService, l logger.Interface) {
	h := &SupplierRouter{service: s, logger: l}
	supplierRoutes := handler.Group("/supplier")
	{
		supplierRoutes.GET("/", h.GetSuppliers)
	}
}

// BasePath /v1/api

// @Summary Get supplier
// @Description Get supplier
// @Tags supplier
// @Accept json
// @Produce json
// @Success 200 {array} top-up-api_internal_schema.SupplierResponse
// @Router /supplier [get]
func (h *SupplierRouter) GetSuppliers(c *gin.Context) {
	suppliers, err := h.service.GetSuppliers(c)
	if err != nil {
		h.logger.Error(errors.New("error getting suppliers"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, mapper.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, mapper.SuccessResponse(suppliers))
}
