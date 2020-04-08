package goee

import (
	"github.com/deific/goee/config"
	"github.com/deific/goee/core"
	"github.com/deific/goee/filter"
	gm "github.com/deific/goee/middleware"
	render2 "github.com/deific/goee/render"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// 定义路由分组
type RouterGroup struct {
	prefix      string             // 路由前缀
	middlewares []core.HandlerFunc // 中间件
	parent      *RouterGroup       // 父组
	engine      *Engine            // 引擎实例
}

// Engine implement the interface of http.ServerHttp
// 定义引擎,最终实现ServerHttp接口
type Engine struct {
	*RouterGroup // 组合group,引擎也具有分组的能力
	// 配置文件
	conf config.Config
	// 路由信息,使用map保存
	router *router
	// 路由分组
	groups []*RouterGroup
	// 全局过滤器
	filters []core.Filter
	// html模板
	renderManager *render2.GoeeRenderManager
}

// New is the constructor of gooee.Engine
func New() *Engine {
	e := &Engine{router: newRouter()}
	// 保存实例
	e.RouterGroup = &RouterGroup{engine: e}
	// 渲染管理器
	e.renderManager = render2.New()

	e.groups = []*RouterGroup{e.RouterGroup}
	return e
}

// default Engine
func Default() *Engine {
	e := New()
	// 加载默认配置
	e.LoadConfig("conf/goee.yaml")
	// 注册默认过滤器
	e.Filter(filter.LogFilter())
	// 注册默认中间件
	e.Use(gm.Logger(), Recovery(), gm.Limit())
	return e
}

// set function map
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.renderManager.SetFuncMap(funcMap)
}

// set function map
func (engine *Engine) LoadHtmlTemplate(pattern string) {
	engine.renderManager.SetHtmlTemplates(template.Must(template.New("").Funcs(engine.renderManager.GetFuncMap()).ParseGlob(pattern)))
}

// 加载配置文件
func (engine *Engine) LoadConfig(filePath string) {
	engine.conf = config.LoadConfig(filePath)
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
func (group *RouterGroup) Use(middleware ...core.HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
}

// 创建static处理
// 使用中间件函数
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) core.HandlerFunc {
	// 绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	// 文件服务器
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *core.Context) {
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
func (e *Engine) Filter(filter ...core.Filter) {
	e.filters = append(e.filters, filter...)
}

// 执行过滤器
func (e Engine) DoFilter(c *core.Context) {
	c.FilterPos++
	if !c.IsStatic && c.FilterPos < len(e.filters) {
		e.filters[c.FilterPos].DoFilter(c, e)
	} else {
		c.Next()
	}
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, isStatic bool, handler core.HandlerFunc) *RouterGroup {
	pattern := group.prefix + comp
	log.Printf("Add Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, isStatic, handler)
	return group
}

// Get define the method of GET to add router
func (group *RouterGroup) GET(patten string, handler core.HandlerFunc) *RouterGroup {
	group.addRoute("GET", patten, false, handler)
	return group
}

// Post define the method of GET to add router
func (group *RouterGroup) POST(patten string, handler core.HandlerFunc) *RouterGroup {
	group.addRoute("POST", patten, false, handler)
	return group
}

// Run define the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	log.Println("启动服务，监听地址：", addr)
	return http.ListenAndServe(addr, engine)
}

// Run define the method to start a http server
func (engine *Engine) Start() (err error) {
	addr := engine.conf.Host + ":" + engine.conf.Port
	printServerInfo(addr)
	return http.ListenAndServe(addr, engine)
}

// print server info
func printServerInfo(addr string) {
	log.Println("************************************")
	log.Println("*             goee " + VERSION + "          *")
	log.Println("************************************")
	log.Println("启动服务，监听地址：", addr)
}

// implement the interface of http.ServerHttp
// 实现http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []core.HandlerFunc
	c := core.NewContext(w, r)
	// 设置过滤器链
	c.FilterChain = engine
	// 传递渲染器
	c.RenderManager = engine.renderManager

	// 识别请求url所在分组，追加分组处理中间件
	for _, group := range engine.groups {
		// 如果请求符合某一组
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 设置中间件
	c.Handlers = middlewares
	// 执行路由
	engine.router.handle(c)
}
