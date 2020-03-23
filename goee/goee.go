package goee

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler user by goee
// 定义框架处理请求的处理器类型
type HandlerFunc func(*Context)

// Engine implement the interface of http.ServerHttp
// 定义引擎,最终实现ServerHttp接口
type Engine struct {
	// 路由信息,使用map保存
	router *router
}

// New is the constructor of gooee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

// Get define the method of GET to add router
func (engine *Engine) GET(patten string, handler HandlerFunc) {
	engine.router.addRouter("GET", patten, handler)
}

// Post define the method of GET to add router
func (engine *Engine) POST(patten string, handler HandlerFunc) {
	engine.router.addRouter("POST", patten, handler)
}

// Run define the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	log.Println("启动服务，监听地址：", addr)
	return http.ListenAndServe(addr, engine)
}

// implement the interface of http.ServerHttp
// 实现http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handle(c)
}
