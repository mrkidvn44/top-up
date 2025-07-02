package controller

import (
	"net/http"
	"strconv"
	"top-up-api/internal/mapper"
	"top-up-api/internal/service"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PurchaseHistoryRouter struct {
	service service.IPurchaseHistoryService
	logger  logger.Interface
	auth    auth.Interface
}

func NewPurchaseHistoryRouter(handler *gin.RouterGroup, s service.IPurchaseHistoryService, l logger.Interface, a auth.Interface) {
	h := &PurchaseHistoryRouter{service: s, logger: l, auth: a}
	handler.GET("/purchase-history/:user_id", h.GetPurchaseHistory)
}

// BasePath /v1/api

// @Summary Get purchase history
// @Description Get purchase history
// @Tags purchase-history
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Security Bearer
// @Success 200 {object} top-up-api_internal_schema.PaginationResponse
// @Router /purchase-history/{user_id} [get]
func (h *PurchaseHistoryRouter) GetPurchaseHistory(c *gin.Context) {
	token, err := h.auth.AuthenticateService(c)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusUnauthorized, mapper.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	userID, err := h.auth.ValidateUserIDFromToken(token, c.Param("user_id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusUnauthorized, mapper.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, mapper.ErrorResponse(http.StatusBadRequest, "Invalid page number", ""))
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, mapper.ErrorResponse(http.StatusBadRequest, "Invalid page size", ""))
		return
	}

	paginatedResponse, err := h.service.GetPurchaseHistoriesByUserIDPaginated(c, userID, page, pageSize)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, mapper.ErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, paginatedResponse)
}
