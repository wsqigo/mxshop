package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mxshop_api/oss_web/middlewares"
	"mxshop_api/oss_web/router"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"success": true,
		})
	})

	r.LoadHTMLFiles("oss_web/templates/index.html")
	// 配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	r.StaticFS("/static", http.Dir(fmt.Sprintf("oss_web/static")))
	// GET: 请求方式；
	r.GET("", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	//配置跨域
	r.Use(middlewares.Cors())

	apiGroup := r.Group("/oss/v1")
	router.InitOssRouter(apiGroup)
	return r
}
