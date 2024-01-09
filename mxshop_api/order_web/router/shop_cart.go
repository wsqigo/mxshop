package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/order_web/api/shop_cart"
)

func InitCartRouter(r *gin.RouterGroup) {
	orderRouter := r.Group("cart")
	//orderRouter := r.Group("cart").Use(middlewares.JWTAuth())
	{
		orderRouter.GET("/list", shop_cart.ListCart)         // 购物车列表
		orderRouter.POST("/add", shop_cart.AddCartItem)      // 添加商品到购物车
		orderRouter.DELETE("/:id", shop_cart.DeleteCartItem) // 删除条目
		orderRouter.PATCH("/:id", shop_cart.UpdateCartItem)  // 修改条目
	}
}
