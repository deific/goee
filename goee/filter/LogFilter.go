package filter

import (
	"goee"
	"log"
)

type logFilter struct {
}

func LogFilter() goee.Filter {
	return logFilter{}
}

func (l logFilter) DoFilter(c *goee.Context, chain goee.Chain) {
	log.Println("执行过滤器前")
	chain.DoFilter(c)
	log.Println("执行过滤器后")
}
