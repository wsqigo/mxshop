package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/order_web/api/order"
	"mxshop_api/order_web/api/pay"
	"mxshop_api/order_web/middlewares"
)

func InitOrderRouter(r *gin.RouterGroup) {
	orderRouter := r.Group("order").Use(middlewares.JWTAuth())
	{
		orderRouter.GET("/list", order.ListOrder)      // 订单列表
		orderRouter.POST("/create", order.CreateOrder) // 创建订单
		orderRouter.GET("/:id", order.GetCartDetail)   // 创建订单
	}
	payRouter := r.Group("pay")
	{
		payRouter.POST("/alipay/notify", pay.Notify)
	}
}
