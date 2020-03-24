package gm

import (
	"goee"
	"log"
	"time"
)

// 记录访问日志
func Logger() goee.HandlerFunc {
	return func(c *goee.Context) {
		// 开始时间
		t := time.Now()
		log.Print("日志中间件执行前")
		c.Next()
		log.Print("日志中间件执行前")
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
