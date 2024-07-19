package web

import (
	"regexp"
	"strings"
)

// 用来支持支持对路由树的操作
// 代表路由树（森林）
type router struct {
	// Beego Gin HTTP method 对应一棵树
	// GET 有一棵树，POST 也有一棵树
	// method => 路由树根节点
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: make(map[string]*node),
	}
}

// addRoute 注册路由
// method 是 HTTP 方法
// path 是路由路径，必须以 / 开头，中间不能有连续的 /
func (r *router) addRoute(method, path string, handler HandlerFunc) {
	if path == "" {
		panic("web: path is empty")
	}
	if path[0] != '/' {
		panic("web: path must begin with '/'")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: path must not end with '/'")
	}
	root, ok := r.trees[method]
	if !ok {
		// 这是一个全新的 HTTP 方法，创建根节点
		root = &node{path: "/", children: make(map[string]*node)}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: path '/' is already registered")
		}
		root.handler = handler
		root.route = path
		return
	}
	path = path[1:]
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			// 不能有连续的 /
			panic("web: path should not have continuous '/'")
		}
		root = root.childOrCreate(seg)
	}
	if root.handler != nil {
		panic("web: path is already registered")
	}
	root.route = path
	root.handler = handler
}

// findRoute 查找路由
// 注意，返回的 node 内部 Handler 不为 nil, 才算是命中路由
// 需要调用者进一步判断
func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}

	mi := &matchInfo{}
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		root, ok = root.childOf(seg)
		if !ok {
			return nil, false
		}
		if root.paramName != "" {
			mi.addValue(root.paramName, seg)
		}
	}

	mi.n = root
	return mi, true
}

// 节点类型
const (
	// 静态节点
	nodeTypeStatic = iota
	// 正则表达式节点
	nodeTypeRegexp
	// 路径参数节点
	nodeTypeParam
	// 通配符节点
	nodeTypeAny
)

type node struct {
	// 节点的路径
	path string

	// 子 path 到子节点的映射
	children map[string]*node
	// 通配符节点
	wildChild *node
	// 正则表达式节点
	regexChild *node
	regExpr    *regexp.Regexp
	// 路径参数节点，跟通配符节点不能同时存在
	paramChild *node
	// 路径参数
	paramName string

	route string

	// 节点的处理函数
	handler HandlerFunc
}

func (n *node) childOrCreate(path string) *node {
	if path == "*" {
		if n.paramChild != nil {
			panic("web: paramChild already exist")
		}
		if n.regexChild != nil {
			panic("web: regexChild already exist")
		}
		if n.wildChild == nil {
			n.wildChild = &node{path: path, children: make(map[string]*node)}
		}
		return n.wildChild
	}

	// 需要区分是路径参数节点还是正则表达式节点
	if path[0] == ':' {
		paramName, expr, isReg := n.parseParam(path)
		if isReg {
			n.childOrCreateReg(path, expr, paramName)
		}
		return n.childOrCreateParam(path, paramName)
	}

	res, ok := n.children[path]
	if !ok {
		res = &node{
			path:     path,
			children: make(map[string]*node),
		}
		n.children[path] = res
	}
	return res
}

// childOf 查找子节点
func (n *node) childOf(seg string) (*node, bool) {
	res, ok := n.children[seg]
	if !ok {
		if n.regexChild != nil {
			if n.regexChild.regExpr.MatchString(seg) {
				return n.regexChild, true
			}
		}
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.wildChild, n.wildChild != nil
	}
	return res, ok
}

// parseParam 用于解析判断是不是正则表达式
// 第一个返回值是参数名
// 第二个返回值是正则表达式
// 第三个返回值是是否为正则表达式
func (n *node) parseParam(path string) (string, string, bool) {
	// 去掉 :
	path = path[1:]

	segs := strings.SplitN(path, "(", 2)
	paramName := segs[0]
	if len(segs) == 2 {
		segs[1] = segs[1][:len(segs[1])-1]
		return paramName, segs[1], true
	}
	return paramName, "", false
}

func (n *node) childOrCreateReg(path string, expr string, paramName string) *node {
	if n.wildChild != nil {
		panic("web: wildChild already exist")
	}
	if n.paramChild != nil {
		panic("web: paramChild already exist")
	}
	if n.regexChild != nil {
		if n.regexChild.regExpr.String() != expr || n.regexChild.paramName != paramName {
			panic("web: regexChild already exist")
		}
	}

	return &node{
		path:      path,
		paramName: paramName,
		regExpr:   regexp.MustCompile(expr),
		children:  make(map[string]*node),
	}
}

func (n *node) childOrCreateParam(path string, paramName string) *node {
	if n.wildChild != nil {
		panic("web: wildChild already exist")
	}
	if n.regexChild != nil {
		panic("web: regexChild already exist")
	}
	if n.paramChild != nil {
		// 判断是否已经有相同参数路由
		if n.paramChild.path != path {
			panic("web: paramChild already exist")
		}
	} else {
		n.paramChild = &node{path: path, paramName: paramName, children: make(map[string]*node)}
	}
	return n.paramChild
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, val string) {
	if m.pathParams == nil {
		m.pathParams = make(map[string]string)
	}
	m.pathParams[key] = val
}
