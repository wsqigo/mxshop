package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user_web/api"
)

func InitBaseRouter(r *gin.RouterGroup) {
	baseRouter := r.Group("base")
	{
		baseRouter.GET("captcha", api.GetCaptcha)
		baseRouter.POST("send_sms", api.SendSms)
	}
}
