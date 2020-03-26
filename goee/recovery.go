package goee

import (
	"fmt"
	"goee/core"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 跟踪信息
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var strb strings.Builder
	strb.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		strb.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return strb.String()
}

// 异常恢复处理
func Recovery() core.HandlerFunc {
	return func(c *core.Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprint("%s", err)
				log.Println("%s\n\n", trace(message))
				c.HTML(http.StatusInternalServerError, "Internal server error")
			}
		}()

		// 继续处理
		c.Next()
	}
}
