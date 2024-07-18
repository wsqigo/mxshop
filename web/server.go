package web

import "net/http"

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
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP 处理请求的入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	// 查找路由，执行代码
	s.serve(ctx)
}

func (s *HTTPServer) serve(ctx *Context) {
	// 先查找路由树
	info, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("404 not found"))
		return
	}
	ctx.PathParams = info.pathParams
	info.n.handler(ctx)
}

// AddRoute 注册路由
func (s *HTTPServer) addRoute(method, path string, handler HandlerFunc) {
	s.router.addRoute(method, path, handler)
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

type HandlerFunc func(ctx *Context)
