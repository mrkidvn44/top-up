package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PurchaseHistoryRouter struct {
	service service.PurchaseHistoryService
	logger  logger.Interface
	auth    auth.Interface
}

func NewPurchaseHistoryRouter(handler *gin.RouterGroup, s service.PurchaseHistoryService, l logger.Interface, a auth.Interface) {
	h := &PurchaseHistoryRouter{service: s, logger: l, auth: a}
	handler.GET("/purchase-history/:id", h.GetPurchaseHistory)
}

// BasePath /v1/api

// @Summary Get purchase history
// @Description Get purchase history
// @Tags purchase-history
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security Bearer
// @Success 200 {array} top-up-api_internal_model.PurchaseHistory
// @Router /purchase-history/{id} [get]
func (h *PurchaseHistoryRouter) GetPurchaseHistory(c *gin.Context) {
	token, err := h.auth.AuthenticateService(c)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	userAuth, err := h.auth.GetUserFromToken(token)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse(http.StatusBadRequest, "Invalid ID", ""))
		return
	}
	if userAuth.ID != uint(idInt) {
		h.logger.Error(fmt.Errorf("userAuth.ID: %v, idInt: %v", userAuth.ID, idInt))
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse(http.StatusUnauthorized, "Unauthorized", ""))
		return
	}

	purchaseHistory, err := h.service.GetPurchaseHistoryByUserID(c, userAuth.ID)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, schema.SuccessResponse(purchaseHistory))
}
