package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/user_web/api"
	"mxshop_api/user_web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		userRouter.GET("list", middlewares.JWTAuth(), middlewares.AdminAuth(), api.GetUserList)
		userRouter.POST("pwd-login", api.PasswordLogin)
		userRouter.POST("register", api.Register)
	}
}
