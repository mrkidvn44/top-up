package controller

import (
	docs "top-up-api/docs"
	"top-up-api/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, services *service.Container) {
	// Swagger
	docs.SwaggerInfo.BasePath = "/v1/api"
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	h := handler.Group("/v1/api")
	{
		NewUserRouter(h, services.UserService, services.Logger, services.Redis, services.Auth, services.Validator)
		NewProviderRouter(h, services.ProviderService, services.Logger)
		NewSkuRouter(h, services.SkuService, services.Logger)
		NewPurchaseHistoryRouter(h, services.PurchaseHistoryService, services.Logger, services.Auth)
		NewOrderRouter(h, services.OrderService, services.Logger, services.Auth, services.Validator)
	}
}