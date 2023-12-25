package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods_web/api/category"
)

func InitCategoryRouter(r *gin.RouterGroup) {
	categoryRouter := r.Group("category")
	{
		categoryRouter.GET("/list", category.ListCategory)            // 商品类别列表页
		categoryRouter.DELETE("/:id", category.DeleteCategory)        // 删除商品分类
		categoryRouter.GET("/detail/:id", category.GetCategoryDetail) // 获取分类详情
		categoryRouter.POST("/create", category.CreateCategory)       // 新建分类
		categoryRouter.PUT("/:id", category.UpdateCategory)           //修改分类信息
	}
}
