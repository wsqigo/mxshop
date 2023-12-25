package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods_web/api/banner"
	"mxshop_api/user_web/middlewares"
)

func InitBannerRouter(r *gin.RouterGroup) {
	bannerRouter := r.Group("banner")
	{
		bannerRouter.GET("/list", banner.ListBanner) // 轮播图列表页
		//bannerRouter.POST("/create", middlewares.JWTAuth(),
		//	middlewares.AdminAuth(), banner.CreateBanner) // 新建轮播图
		bannerRouter.POST("/create", banner.CreateBanner)
		//bannerRouter.DELETE("/:id", middlewares.JWTAuth(),
		//	middlewares.AdminAuth(), banner.DeleteBanner) // 删除轮播图
		bannerRouter.DELETE("/:id", banner.DeleteBanner) // 删除轮播图
		bannerRouter.PUT("/:id", middlewares.JWTAuth(),
			middlewares.AdminAuth(), banner.UpdateBanner)
	}
}
