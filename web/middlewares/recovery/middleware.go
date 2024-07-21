package recovery

import (
	"awesomeProject/web"
	"net/http"
)

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *web.Context)
}

func NewMiddlewareBuilder(errMsg string, logFunc func(ctx *web.Context)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		StatusCode: http.StatusInternalServerError,
		ErrMsg:     errMsg,
		LogFunc:    logFunc,
	}
}

func (b *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = b.StatusCode
					ctx.RespData = []byte(b.ErrMsg)
					// 万一 LogFunc 也 panic，那我们也无能为力了
					b.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}
