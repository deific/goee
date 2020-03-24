package gm

import (
	"goee"
	"log"
)

func Limit() goee.HandlerFunc {
	return func(c *goee.Context) {
		log.Println("限流中间件执行前")
		c.Next()
		log.Println("限流中间件执行后")
	}
}
