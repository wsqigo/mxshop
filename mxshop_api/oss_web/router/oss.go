package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/oss_web/api"
)

func InitOssRouter(r *gin.RouterGroup) {
	ossRouter := r.Group("oss")
	{
		ossRouter.GET("/token", api.GetToken)
		ossRouter.POST("/callback", api.HandlerRequest)
	}
}
