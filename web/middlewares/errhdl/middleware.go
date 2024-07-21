package errhdl

import "awesomeProject/web"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		// 这里可以非常大方，因为在预计中用户会关心的错误码不可能超过 64
		resp: make(map[int][]byte, 64),
	}
}

// RegisterError 注册错误码和响应内容
// 这个错误数据可以是一个字符串，也可以是一个页面
func (b *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	b.resp[code] = resp
	return b
}

func (b *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) {
			next(ctx)
			if resp, ok := b.resp[ctx.RespStatusCode]; ok {
				ctx.RespData = resp
			}
		}
	}
}
