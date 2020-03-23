package main

import (
	"goee"
	"net/http"
)

// 启动入口
func main() {
	// 实例化
	g := goee.New()
	// 注册路由
	g.GET("/", func(c *goee.Context) {
		c.HTML(http.StatusOK, "<h1>Welcome to Goee!</h1> you request path = "+c.Path)
	})

	g.GET("/hello/:name", func(c *goee.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})

	g.GET("/json", func(c *goee.Context) {
		c.JSON(http.StatusOK, goee.HMap{"code": 200, "msg": "ok", "success": true})
	})

	// 启动服务
	g.Run(":9000")
}
