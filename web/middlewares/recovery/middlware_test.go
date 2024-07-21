package recovery

import (
	"awesomeProject/web"
	"log"
	"testing"
)

func TestMiddleware(t *testing.T) {
	m := NewMiddlewareBuilder("internal server error", func(ctx *web.Context) {
		log.Println(ctx.Req.URL.Path, "internal server error")
	}).Build()

	s := web.NewHTTPServer(web.ServerWithMiddleware(m))
	s.Get("/user", func(ctx *web.Context) {
		panic("something went wrong")
	})

	s.Start(":8081")
}
