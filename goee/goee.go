package goee

import (
	"log"
	"net/http"
	"path"
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

// 创建static处理
// 使用中间件函数
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// 绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	// 文件服务器
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filePath")
		// Check if file exist and/or if we have permission to access
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		// 返回文件服务
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 设置静态文件目录
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.addRoute("GET", urlPattern, true, handler)
}

// 设置全局过滤器
func (e *Engine) Filter(filter ...Filter) {
	e.filters = append(e.filters, filter...)
}

// 执行过滤器
func (e Engine) DoFilter(c *Context) {
	c.filterPos++
	if !c.isStatic && c.filterPos < len(e.filters) {
		e.filters[c.filterPos].DoFilter(c, e)
	} else {
		c.Next()
	}
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, isStatic bool, handler HandlerFunc) *RouterGroup {
	pattern := group.prefix + comp
	log.Printf("Add Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, isStatic, handler)
	return group
}

// Get define the method of GET to add router
func (group *RouterGroup) GET(patten string, handler HandlerFunc) *RouterGroup {
	group.addRoute("GET", patten, false, handler)
	return group
}

// Post define the method of GET to add router
func (group *RouterGroup) POST(patten string, handler HandlerFunc) *RouterGroup {
	group.addRoute("POST", patten, false, handler)
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
	c := newContext(w, r, engine)

	for _, group := range engine.groups {
		// 如果请求符合某一组
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c.handlers = middlewares
	// 执行路由
	engine.router.handle(c)
}
