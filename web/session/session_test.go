package session

import (
	"awesomeProject/web"
	"net/http"
	"testing"
)

func TestSession(t *testing.T) {
	var m Manager
	s := web.NewHTTPServer(web.ServerWithMiddleware(func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) {
			// 执行校验
			if ctx.Req.URL.Path == "/login" {
				_, err := m.GetSession(ctx)
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					ctx.RespData = []byte("please login")
					return
				}

				err = m.RefreshSession(ctx)
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					ctx.RespData = []byte("please login")
					return
				}
			}
			next(ctx)
		}
	}))

	s.Post("/login", func(ctx *web.Context) {
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("init session error")
			return
		}
		// 然后根据自己的需要设置
		err = sess.Set(ctx.Req.Context(), "mykey", "some value")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("set session error")
			return
		}
	})

	s.Start(":8080")
}
