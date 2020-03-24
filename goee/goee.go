package goee

import (
	"log"
	"net/http"
	"strings"
)

// 定义路由分组
type RouterGroup struct {
	prefix      string        // 路由前缀
	middlewares []HandlerFunc // 中间件
	parent      *RouterGroup  // 父组
	engine      *Engine       // 引擎实例
}

// HandlerFunc defines the request handler user by goee
// 定义框架处理请求的处理器类型
type HandlerFunc func(*Context)

// Engine implement the interface of http.ServerHttp
// 定义引擎,最终实现ServerHttp接口
type Engine struct {
	*RouterGroup // 组合group,引擎也具有分组的能力
	// 路由信息,使用map保存
	router *router
	// 路由分组
	groups []*RouterGroup
	// 全局过滤器
	filters []Filter
}

// New is the constructor of gooee.Engine
func New() *Engine {
	e := &Engine{router: newRouter()}
	// 保存实例
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}
	return e
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	e := group.engine
	newGroup := &RouterGroup{
		prefix: prefix,
		parent: group,
		engine: e,
	}
	e.groups = append(e.groups, newGroup)
	return newGroup
}

// 使用中间件函数
func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
}

// 设置全局过滤器
func (e *Engine) Filter(filter ...Filter) {
	e.filters = append(e.filters, filter...)
}

// 执行过滤器
func (e Engine) DoFilter(c *Context) {
	c.filterPos++
	if c.filterPos < len(e.filters) {
		e.filters[c.filterPos].DoFilter(c, e)
	} else {
		c.Next()
	}
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) *RouterGroup {
	pattern := group.prefix + comp
	log.Printf("Add Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
	return group
}

// Get define the method of GET to add router
func (group *RouterGroup) GET(patten string, handler HandlerFunc) *RouterGroup {
	group.addRoute("GET", patten, handler)
	return group
}

// Post define the method of GET to add router
func (group *RouterGroup) POST(patten string, handler HandlerFunc) *RouterGroup {
	group.addRoute("POST", patten, handler)
	return group
}

// Run define the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	log.Println("启动服务，监听地址：", addr)
	return http.ListenAndServe(addr, engine)
}

// implement the interface of http.ServerHttp
// 实现http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc

	for _, group := range engine.groups {
		// 如果请求符合某一组
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, r, engine)
	c.handlers = middlewares
	// 执行路由
	engine.router.handle(c)

}
