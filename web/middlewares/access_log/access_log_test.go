package access_log

import (
	"awesomeProject/web"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccessLogMiddleware_SuccessfulRequest(t *testing.T) {
	// Mock log function
	logFunc := func(log string) {
		fmt.Println(log)
	}

	// Create AccessLogMiddlewareBuilder with mocked log function
	builder := NewAccessLogMiddlewareBuilder().
		WithLogFunc(logFunc)
	middleware := builder.Build()

	// Create a mock handler
	handler := func(ctx *web.Context) {
		fmt.Println("业务代码")
	}

	// Create a mock Context
	ctx := &web.Context{
		Req: httptest.NewRequest(http.MethodGet, "/test", nil),
	}

	// Call the middleware
	middleware(handler)(ctx)
}
