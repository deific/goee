package goee

import "net/http"

// 上下文

type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    http.Request

	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

// 创建新的上下文
func newContext(w http.ResponseWriter, r http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        r,
		Path:       r.URL.Path,
		Method:     r.Method,
		StatusCode: 200,
	}
}

// 获取form表单数据
