package controller

import (
	docs "cashier-api/docs"
	"cashier-api/internal/service"

	orderpb "cashier-api/proto/order"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, services *service.Container, grpcOrderService orderpb.OrderServiceClient) {
	// Swagger
	docs.SwaggerInfo.BasePath = "/v1/api"
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	h := handler.Group("/v1/api")
	{
		NewOrderRouter(h, services.OrderService, services.Logger, grpcOrderService)
	}
}
