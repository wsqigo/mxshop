package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,x-token")
		ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PATCH,PUT")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Header,Content-Type")
		ctx.Header("Access-Control-All-Credentials", "true")
		if method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
