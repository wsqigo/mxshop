package web

import (
	"log"
	"net/http"
)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 监听地址。如果只指定端口，可以使用 ":port"
	// 正常 "ip:port"
	Start(addr string) error
	// AddRoute 注册路由
	// method 是 HTTP 方法
	// path 是路由路径，必须以 / 开头
	addRoute(method, path string, handler HandlerFunc)
}

type HTTPServer struct {
	// addr string 创建的时候传递，还是 Start 接收。这个都是可以的
	router

	mdls []Middleware // 中间件

	tplEngine TemplateEngine
}

// HTTPServerOption 模式
type HTTPServerOption func(*HTTPServer)

// ServerWithMiddleware 添加中间件
func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(s *HTTPServer) {
		s.mdls = append(s.mdls, mdls...)
	}
}

// ServerWithTemplateEngine 设置模板引擎
func ServerWithTemplateEngine(engine TemplateEngine) HTTPServerOption {
	return func(s *HTTPServer) {
		s.tplEngine = engine
	}
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServeHTTP 处理请求的入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:            r,
		Resp:           w,
		templateEngine: s.tplEngine,
	}

	root := s.serve
	// 洋葱模式，将中间件组合成链
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}

	var m Middleware = func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	root = m(root)
	// 查找路由，执行代码
	root(ctx)
}

func (s *HTTPServer) serve(ctx *Context) {
	// 先查找路由树
	info, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("not found")
		return
	}
	ctx.PathParams = info.pathParams
	ctx.MatchedRoute = info.n.route
	info.n.handler(ctx)
}

// AddRoute 注册路由
func (s *HTTPServer) addRoute(method, path string, handler HandlerFunc, mdls ...Middleware) {
	s.router.addRoute(method, path, handler, mdls...)
}

func (s *HTTPServer) Post(path string, handler HandlerFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandlerFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

// Use 在路由树中注册中间件
func (s *HTTPServer) Use(method, path string, mdls ...Middleware) {
	s.addRoute(method, path, nil, mdls...)
}

func (s *HTTPServer) flashResp(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Println("flashResp error:", n, err)
	}
}

type HandlerFunc func(ctx *Context)
