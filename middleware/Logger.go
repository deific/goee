package gm

import (
	"github.com/deific/goee/core"
	"log"
	"time"
)

// 记录访问日志
func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		// 开始时间
		t := time.Now()
		log.Print("日志中间件执行前")
		c.Next()
		log.Print("日志中间件执行前")
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
