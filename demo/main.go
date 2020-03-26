package main

import (
	"github.com/deific/goee"
	"github.com/deific/goee/core"
	"log"
	"net/http"
	"time"
)

// 启动入口
func main() {
	// 实例化
	g := goee.Default()

	// 注册静态文件处理
	g.Static("/static", "E:/temp/static")
	// 加载模板
	g.LoadHtmlTemplate("templates/*")

	// 注册路由分组
	v1 := g.Group("/v1")
	v1.GET("/hello", func(c *core.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.GetParam("name")+"</h1>")
	}).GET("/hello/:name", func(c *core.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})

	// 注册路由，想到根分组
	g.GET("/", func(c *core.Context) {
		c.HTML(http.StatusOK, "<h1>Welcome to Goee!</h1> you request path = "+c.Path)
	})
	g.GET("/hello/:name", func(c *core.Context) {
		log.Println("执行实际处理函数")
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})
	g.GET("/json", func(c *core.Context) {
		c.JSON(http.StatusOK, core.HMap{"code": 200, "msg": "成功", "success": true})
	})
	// 使用html模板
	g.GET("/html", func(c *core.Context) {
		c.HtmlTemplate(http.StatusOK, "index.tpl", core.HMap{"title": "张三", "now": time.Now()})
	})

	// 异常路由
	g.GET("/panic", func(c *core.Context) {
		name := []string{"goee"}
		c.JSON(http.StatusOK, name[10])
	})

	// 启动服务
	g.Run(":9000")
}
