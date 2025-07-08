package controller

import (
	"net/http"
	docs "top-up-api/docs"
	grpcClient "top-up-api/internal/grpc/client"
	"top-up-api/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, services *service.Container, grpcClients *grpcClient.GRPCServiceClient) {
	// Health check endpoint
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

	// Swagger
	docs.SwaggerInfo.BasePath = "/v1/api"
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	h := handler.Group("/v1/api")
	{
		NewProviderRouter(h, services.ProviderService, services.Logger)
		NewSkuRouter(h, services.SkuService, services.Logger)
		NewPurchaseHistoryRouter(h, services.PurchaseHistoryService, grpcClients.AuthGRPCClient, services.Logger)
		NewOrderRouter(h, services.OrderService, grpcClients.AuthGRPCClient, services.Logger, services.Validator)
	}
}
