package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/userop_web/api/user_fav"
	"mxshop_api/userop_web/middlewares"
)

func InitUserFavRouter(r *gin.RouterGroup) {
	userFavRouter := r.Group("userfav").Use(middlewares.JWTAuth())
	{
		userFavRouter.DELETE("/:id", user_fav.DeleteUserFav)  // 删除收藏记录
		userFavRouter.GET("/:id", user_fav.GetUserFavDetail)  // 获取收藏记录
		userFavRouter.POST("/create", user_fav.CreateUserFav) //新建收藏记录
		userFavRouter.GET("/list", user_fav.ListUserFav)      //获取当前用户的收藏
	}
}
