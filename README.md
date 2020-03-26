# Goee web framework

Goee 是基于GO语言开发的web开发框架，参考和借鉴了Gin的设计思路。对Gin进行简化和优化调整，是一个更轻量更简洁高性能的web开发框架。

# Installation

To install Goee package, you need to install Go and set your Go workspace first.

The first need Go installed (version 1.13+ is required), then you can use the below Go command to install Goee.
```go
$ go get -u github.com/gin-gonic/gin
```
Import it in your code:
```go
import "github.com/deific/goee/tree/master/goee"
```

# Quick start
 ```go
# assume the following codes in example.go file
$ cat example.go
```

```go
package main

import (
	"github.com/deific/goee/tree/master/goee"
	"github.com/deific/goee/tree/master/goee/core"
	"log"
	"net/http"
	"time"
)
func main() {
	// 实例化
    g := goee.Default()

    // 注册路由分组
    v1 := g.Group("/v1")
    v1.GET("/hello", func(c *core.Context) {
        c.HTML(http.StatusOK, "<h1>hello "+c.GetParam("name")+"</h1>")
    }).GET("/hello/:name", func(c *core.Context) {
        c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
    })
}
```


