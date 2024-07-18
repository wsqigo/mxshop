package web

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func readBodyOnce(w http.ResponseWriter, r *http.Request) {
	r.URL.Query()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body error: %v", err)
		// 记住要返回，不然就还会执行后面的代码
		return
	}
	// 类型转换，将 []byte 转为 string
	fmt.Fprintf(w, "read body: %s", body)

	// 尝试再次读取，啥也读不到，但是也不会报错
	body, err = io.ReadAll(r.Body)
	if err != nil {
		// 不会进来这里
		fmt.Fprintf(w, "read the data one more time error: %v", err)
		return
	}
	fmt.Fprintf(w, "read the data one more time: [%s] and read data length: %d", body, len(body))
}

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/user/:id", func(ctx *Context) {
		id, err := ctx.PathValue("id").ToInt64()
		if err != nil {
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(id)
	})
	s.Start(":8081")
}
