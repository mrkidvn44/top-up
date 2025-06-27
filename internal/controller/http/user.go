package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"top-up-api/internal/schema"
	service "top-up-api/internal/service"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/logger"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRouter struct {
	service   service.UserService
	logger    logger.Interface
	redis     redis.Interface
	auth      auth.Interface
	validator validator.Interface
}

func NewUserRouter(handler *gin.RouterGroup, service service.UserService, l logger.Interface, r redis.Interface, a auth.Interface, v validator.Interface) {
	h := &UserRouter{service: service, logger: l, redis: r, auth: a, validator: v}
	UserRoutes := handler.Group("/user")
	{
		UserRoutes.POST("/login", h.Login)
		UserRoutes.GET("/:id", h.GetUserByID)
		UserRoutes.POST("/create", h.CreateUser)
	}
}

// @BasePath /v1/api

// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} top-up-api_internal_schema.UserProfileResponse
// @Failure 401 {object} gin.H
// @Router /user/{id} [get]
// @Security Bearer
func (h *UserRouter) GetUserByID(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if userAuth.ID != uint(idInt) {
		h.logger.Error(fmt.Errorf("userAuth.ID: %v, idInt: %v", userAuth.ID, idInt))
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse(http.StatusUnauthorized, "Unauthorized", ""))
		return
	}

	user, err := h.service.GetUserByID(c, uint(idInt))
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, schema.SuccessResponse(user))
}

// @Summary Login
// @Description Login
// @Tags user
// @Accept json
// @Produce json
// @Param user body top-up-api_internal_schema.UserLoginRequest true "User login request"
// @Success 200 {object} top-up-api_internal_schema.UserLoginDetail
// @Router /user/login [post]
func (h *UserRouter) Login(c *gin.Context) {
	var user schema.UserLoginDetail
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Validate(user); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse(http.StatusBadRequest, "Invalid request", err.Error()))
		return
	}

	userDetail, err := h.service.Login(c, user.PhoneNumber, user.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Error(err)
			c.JSON(http.StatusNotFound, schema.ErrorResponse(http.StatusNotFound, "User not found", ""))
			return
		}
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	token, err := h.auth.CreateToken(*userDetail)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, schema.SuccessResponse(gin.H{"token": token, "message": "Login successful"}))
}

// @Summary Create user
// @Description Create user
// @Tags user
// @Accept json
// @Produce json
// @Param user body top-up-api_internal_schema.UserCreateRequest true "User create details"
// @Success 200 string message "User created successfully"
// @Router /user/create [post]
func (h *UserRouter) CreateUser(c *gin.Context) {
	var user schema.UserCreateRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Validate(user); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse(http.StatusBadRequest, "Invalid request", err.Error()))
		return
	}

	err := h.service.CreateUser(c, user)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
