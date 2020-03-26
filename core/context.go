package core

import (
	"encoding/json"
	"fmt"
	render2 "goee/render"
	"net/http"
)

// 针对map[string]interface定义别名
type HMap map[string]interface{}

// 上下文
type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path     string
	Method   string
	Params   map[string]string
	IsStatic bool //是否静态文件请求
	// response info
	StatusCode int
	// middleware
	Handlers []HandlerFunc
	index    int
	// 过滤器执行索引
	FilterChain Chain
	FilterPos   int
	// 引擎实例
	Render render2.Render
}

// 创建新的上下文
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        r,
		Path:       r.URL.Path,
		Method:     r.Method,
		StatusCode: 200,
		index:      -1,
		FilterPos:  -1,
	}
}

// 中间调用链
func (c *Context) Next() {
	c.index++
	c.Handlers[c.index](c)
}

// 获取url参数
func (c Context) Param(key string) string {
	return c.Params[key]
}

// 获取form表单数据
func (c Context) GetParam(key string) string {
	return c.Req.FormValue(key)
}

// 获取请求参数
func (c Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置响应码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置header
func (c *Context) setHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 文本渲染
func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprint(format, values)))
}

// json渲染
func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// html渲染
func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// html模板渲染
func (c *Context) HtmlTemplate(code int, name string, data interface{}) {
	c.setHeader("Content-Type", "text/html;charset=utf8")
	c.Status(code)
	// 使用模板渲染
	if err := c.Render.GetHtmlTemplates().ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Status(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
	}
}
