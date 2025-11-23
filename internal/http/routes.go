package http

import (
	"insider/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Routers struct {
	healthHandler   *HealthHandler
	messagesHandler *MessagesHandler
	senderHandler   *SenderHandler
}

func NewRouters(hh *HealthHandler, mh *MessagesHandler, sh *SenderHandler) *Routers {
	return &Routers{
		healthHandler:   hh,
		messagesHandler: mh,
		senderHandler:   sh,
	}
}

func (r *Routers) Register(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	health := api.Group("/health")
	health.GET("", r.healthHandler.Health)

	messages := api.Group("/messages")
	messages.GET("", r.messagesHandler.ListMessages)

	sender := api.Group("/sender")
	sender.POST("/toggle", r.senderHandler.ToggleSender)
}
