package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user_web/middlewares"
	"mxshop_api/user_web/router"
	"net/http"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 配置跨域
	r.Use(middlewares.Cors())
	apiGroup := r.Group("/u/v1")
	router.InitUserRouter(apiGroup)
	router.InitBaseRouter(apiGroup)

	return r
}
