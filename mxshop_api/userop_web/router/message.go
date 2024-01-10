package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/userop_web/api/message"
	"mxshop_api/userop_web/middlewares"
)

func InitMessageRouter(r *gin.RouterGroup) {
	messageRouter := r.Group("message").Use(middlewares.JWTAuth())

	{
		messageRouter.GET("list", message.ListMessage)
		messageRouter.POST("create", message.CreateMessage)
	}
}
