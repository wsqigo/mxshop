package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods_web/middlewares"
	"mxshop_api/goods_web/router"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 配置跨域
	r.Use(middlewares.Cors())

	apiGroup := r.Group("/g/v1")
	router.InitGoodsRouter(apiGroup)
	router.InitCategoryRouter(apiGroup)
	router.InitBannerRouter(apiGroup)
	router.InitBrandRouter(apiGroup)
	return r
}
