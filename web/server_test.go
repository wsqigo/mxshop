package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareChain_Success(t *testing.T) {
	// Create a new HTTPServer instance with mock middleware
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				fmt.Println("middleware1 start")
				next(ctx)
				fmt.Println("middleware1 end")
			}
		},
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				fmt.Println("middleware2 start")
				next(ctx)
				fmt.Println("middleware2 end")
			}
		},
	}

	// Add route
	server.Get("/test", func(ctx *Context) {
		fmt.Println("test")
	})

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Call the server
	server.ServeHTTP(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %d", resp.StatusCode)
	}
}
