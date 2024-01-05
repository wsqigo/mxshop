package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/order_web/api/goods"
	"mxshop_api/order_web/middlewares"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	goodsRouter := r.Group("goods")
	{
		goodsRouter.GET("list", goods.ListGoods) /// 商品列表
		//goodsRouter.POST("create", middlewares.JWTAuth(),
		//	middlewares.AdminAuth(), goods.CreateGoods) // 需要管理员权限
		goodsRouter.POST("create", goods.CreateGoods)
		goodsRouter.GET("/:id", goods.Detail) // 获取商品的详情
		//goodsRouter.DELETE("/:id", middlewares.JWTAuth(),
		//	middlewares.AdminAuth(), goods.DeleteGoods) // 删除商品
		goodsRouter.DELETE("/:id", goods.DeleteGoods)
		goodsRouter.GET("/:id/stock", goods.GetGoodsStocks) // 获取商品库存

		goodsRouter.PUT("/:id", middlewares.JWTAuth(),
			middlewares.AdminAuth(), goods.UpdateGoods)
		goodsRouter.PATCH("/:id", middlewares.JWTAuth(),
			middlewares.AdminAuth(), goods.UpdateStatus)
	}
}
