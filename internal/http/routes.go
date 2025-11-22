package http

import (
	"insider/docs"
	"insider/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Routers struct {
	messagesHandler *handler.MessagesHandler
	senderHandler   *handler.SenderHandler
}

func NewRouters(mh *handler.MessagesHandler, sh *handler.SenderHandler) *Routers {
	return &Routers{
		messagesHandler: mh,
		senderHandler:   sh,
	}
}

func (r *Routers) Register(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	health := api.Group("/health")
	health.GET("", handler.Health)

	messages := api.Group("/messages")
	messages.GET("", r.messagesHandler.ListMessages)

	sender := api.Group("/sender")
	sender.POST("/start", r.senderHandler.StartSender)
	sender.POST("/stop", r.senderHandler.StopSender)
}
