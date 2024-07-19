package main

import (
	"fmt"
	"math"
)

// gin 洋葱模式

const abortIndex = math.MaxInt8 / 2

type Context struct {
	handlers []func(ctx *Context)
	index    int8
}

func (ctx *Context) Use(f func(ctx *Context)) {
	ctx.handlers = append(ctx.handlers, f)
}

func (ctx *Context) GET(path string, f func(ctx *Context)) {
	// 逻辑代码...
	ctx.handlers = append(ctx.handlers, f)
}

func (ctx *Context) Next() {
	ctx.index++
	if ctx.index < int8(len(ctx.handlers)) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) Abort() {
	ctx.index = abortIndex
}

func (ctx *Context) Run() {
	ctx.handlers[0](ctx)
}

func Middleware1() func(ctx *Context) {
	return func(ctx *Context) {
		fmt.Println("mid1")
		ctx.Next()
		fmt.Println("mid1 end")
	}
}

func Middleware2() func(ctx *Context) {
	return func(ctx *Context) {
		fmt.Println("mid2")
		ctx.Next()
		fmt.Println("mid2 end")
	}
}

func Middleware3() func(ctx *Context) {
	return func(ctx *Context) {
		fmt.Println("mid3")
		ctx.Next()
		fmt.Println("mid3 end")
	}
}

func main() {
	ctx := &Context{}
	ctx.Use(Middleware1())
	ctx.Use(Middleware2())
	ctx.Use(Middleware3())
	ctx.GET("/", func(ctx *Context) {
		fmt.Println("业务逻辑")
	})
	ctx.Run()
}

// 集中式
// func (ctx *Context) Run() {
// 	for i := 0; i < len(middleware); i++ {
// 		ctx.handlers[i](ctx)
// 	}
// }
