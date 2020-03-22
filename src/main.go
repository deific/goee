package main

import (
	"fmt"
	"goee"
	"net/http"
)

// 启动入口
func main() {
	// 实例化
	goee := goee.New()
	// 注册路由
	goee.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Welcome to Goee! you request path = %s", req.URL.Path)
	})
	// 启动服务
	goee.Run(":9000")
}
