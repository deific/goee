package core

// 定义过滤器接口
type Filter interface {
	DoFilter(c *Context, fc Chain)
}

// 定义过滤器链
type Chain interface {
	DoFilter(c *Context)
}
