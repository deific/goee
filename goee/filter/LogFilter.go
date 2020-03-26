package filter

import (
	"goee/core"
	"log"
)

type logFilter struct {
}

func LogFilter() core.Filter {
	return logFilter{}
}

func (l logFilter) DoFilter(c *core.Context, chain core.Chain) {
	log.Println("执行过滤器前")
	chain.DoFilter(c)
	log.Println("执行过滤器后")
}
