package main

import (
	"goee"
	"goee/filter"
	gm "goee/middleware"
	"log"
	"net/http"
)

// 启动入口
func main() {
	// 实例化
	g := goee.New()

	// 注册过滤器
	g.Filter(filter.LogFilter())
	// 注册中间件
	g.Use(gm.Logger(), gm.Limit())

	// 注册静态文件处理
	g.Static("/static", "E:/temp/static")

	// 注册路由分组
	v1 := g.Group("/v1")
	v1.GET("/hello", func(c *goee.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.GetParam("name")+"</h1>")
	}).GET("/hello/:name", func(c *goee.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})

	// 注册路由，想到根分组
	g.GET("/", func(c *goee.Context) {
		c.HTML(http.StatusOK, "<h1>Welcome to Goee!</h1> you request path = "+c.Path)
	})
	g.GET("/hello/:name", func(c *goee.Context) {
		log.Println("执行实际处理函数")
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})
	g.GET("/json", func(c *goee.Context) {
		c.JSON(http.StatusOK, goee.HMap{"code": 200, "msg": "ok", "success": true})
	})

	// 启动服务
	g.Run(":9000")
}
