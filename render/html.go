package render

import (
	"net/http"
)

// html渲染器
type HtmlRender struct {
}

// 渲染
func (htmlRender HtmlRender) Render(write http.ResponseWriter, result Result) {
	write.Header().Set("Content-Type", "text/html")
	write.WriteHeader(result.code)
	write.Write([]byte(result.data.(string)))
}
