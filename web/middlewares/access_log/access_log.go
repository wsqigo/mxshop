package access_log

import (
	"awesomeProject/web"
	"encoding/json"
)

// MiddlewareBuilder 用来构建 Middleware
type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

// NewMiddlewareBuilder 工厂模式
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

// WithLogFunc 设置 LogFunc
func (b *MiddlewareBuilder) WithLogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

// Build 构建 AccessLogMiddleware
func (b *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) {
			defer func() {
				l := &accessLog{
					Host:   ctx.Req.Host,
					Route:  ctx.MatchedRoute,
					Method: ctx.Req.Method,
					Path:   ctx.Req.URL.Path,
				}
				bs, _ := json.Marshal(l)
				b.logFunc(string(bs))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host   string
	Route  string
	Method string
	Path   string
	Status int
}
