package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/order_web/middlewares"
	"mxshop_api/order_web/router"
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

	apiGroup := r.Group("/o/v1")
	router.InitOrderRouter(apiGroup)
	router.InitCartRouter(apiGroup)
	return r
}
