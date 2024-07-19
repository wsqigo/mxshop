package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	// PathParams 用来存储路由参数
	PathParams map[string]string
	// cacheQueryValues 用来存储查询参数
	cacheQueryValues url.Values
	// MatchedRoute 用来存储匹配的路由
	MatchedRoute string
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseInt(s.val, 10, 64)
}

//---------------------------------处理输入参数---------------------------------

func (c *Context) BindJson(val any) error {
	if c.Req.Body == nil {
		return errors.New("request body is nil")
	}
	// c.Req.Body 是 io.ReadCloser 接口
	// 不能用 json.Unmarshal 直接解析
	// bs, _ := io.ReadAll(c.Req.Body)
	// err = json.Unmarshal(bs, val)
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

// FormValue 获取表单参数
func (c *Context) FormValue(key string) StringValue {
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{err: err}
	}

	vals, ok := c.Req.Form[key]
	if !ok || len(vals) == 0 {
		return StringValue{err: errors.New("key not found")}
	}

	return StringValue{val: vals[0]}
}

// QueryValue 获取查询参数
func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}

	vals, ok := c.cacheQueryValues[key]
	if !ok || len(vals) == 0 {
		return StringValue{err: errors.New("key not found")}
	}
	return StringValue{val: vals[0]}
}

// PathValue 获取路径参数
func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("key not found")}
	}

	return StringValue{val: val}
}

//---------------------------------输出响应--------------------------------

func (c *Context) RespJson(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	// 设置响应状态
	c.Resp.WriteHeader(code)
	// 设置响应头
	c.Resp.Header().Set("Content-Type", "application/json")
	_, err = c.Resp.Write(bs)

	return err
}

func (c *Context) RespJsonOK(val any) error {
	return c.RespJson(http.StatusOK, val)
}

func (c *Context) String(code int, val string) error {
	// 设置响应状态
	c.Resp.WriteHeader(code)
	// 设置响应头
	c.Resp.Header().Set("Content-Type", "text/plain")
	_, err := c.Resp.Write([]byte(val))

	return err
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}
