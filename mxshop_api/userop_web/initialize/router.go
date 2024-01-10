package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/userop_web/middlewares"
	"mxshop_api/userop_web/router"
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

	apiGroup := r.Group("/up/v1")
	router.InitAddressRouter(apiGroup)
	router.InitMessageRouter(apiGroup)
	router.InitUserFavRouter(apiGroup)
	return r
}
