package goee

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc defines the request handler user by goee
// 定义框架处理请求的处理器类型
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of http.ServerHttp
// 定义引擎,最终实现ServerHttp接口
type Engine struct {
	// 路由信息,使用map保存
	router map[string]HandlerFunc
}

// New is the constructor of gooee.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由
func (engine *Engine) addRouter(method string, patten string, handler HandlerFunc) {
	key := method + "-" + patten
	engine.router[key] = handler
}

// Get define the method of GET to add router
func (engine *Engine) GET(patten string, handler HandlerFunc) {
	engine.addRouter("GET", patten, handler)
}

// Post define the method of GET to add router
func (engine *Engine) POST(patten string, handler HandlerFunc) {
	engine.addRouter("POST", patten, handler)
}

// Run define the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	log.Println("启动服务，监听地址：", addr)
	return http.ListenAndServe(addr, engine)
}

// implement the interface of http.ServerHttp
// 实现http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 解析请求
	key := r.Method + "-" + r.URL.Path
	if handler := engine.router[key]; handler != nil {
		handler(w, r)
	} else {
		fmt.Fprint(w, "404 not found:%s\n", r.URL)
	}
}
