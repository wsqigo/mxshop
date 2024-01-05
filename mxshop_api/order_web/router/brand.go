package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/order_web/api/brand"
)

// 1. 商品的api接口开发完成
// 2. 图片的坑
func InitBrandRouter(r *gin.RouterGroup) {
	brandRouter := r.Group("brand")
	{
		brandRouter.GET("/list", brand.GetBrandList)
		brandRouter.POST("/create", brand.CreateBrand)
		brandRouter.DELETE("/:id", brand.DeleteBrand)
		brandRouter.PUT("/:id", brand.UpdateBrand)
	}

	categoryBrandRouter := r.Group("category-brand")
	{
		categoryBrandRouter.GET("/list-brand", brand.CategoryBrandList) // 类别品牌列表页
		categoryBrandRouter.DELETE("/:id", brand.DeleteCategoryBrand)   // 删除类别品牌
		categoryBrandRouter.POST("/create", brand.CategoryBrandList)    // 新建类别品牌
		categoryBrandRouter.PUT("/:id", brand.UpdateCategoryBrand)      // 修改类别品牌
		categoryBrandRouter.GET("/:id", brand.GetCategoryBrandList)     // 获取分类的品牌
	}
}
