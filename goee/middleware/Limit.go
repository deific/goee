package gm

import (
	"goee/core"
	"log"
)

func Limit() core.HandlerFunc {
	return func(c *core.Context) {
		log.Println("限流中间件执行前")
		c.Next()
		log.Println("限流中间件执行后")
	}
}
