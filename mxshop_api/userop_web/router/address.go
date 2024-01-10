package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/userop_web/api/address"
	"mxshop_api/userop_web/middlewares"
)

func InitAddressRouter(r *gin.RouterGroup) {
	addressRouter := r.Group("address").Use(middlewares.JWTAuth())
	{
		addressRouter.GET("/list", address.ListAddress)      // 轮播图列表页
		addressRouter.DELETE("/:id", address.DeleteAddress)  // 删除轮播图
		addressRouter.POST("/create", address.CreateAddress) //新建轮播图
		addressRouter.PUT("/:id", address.UpdateAddress)     //修改轮播图信息
	}
}
