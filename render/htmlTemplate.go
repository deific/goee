package render

import (
	"html/template"
	"net/http"
)

// html渲染器
type HtmlTemplateRender struct {
	// html模板
	HtmlTemplates *template.Template
	FuncMap       template.FuncMap
}

// 模板结果
type HtmlTplResult struct {
	Result
	templateName string
}

// 渲染
func (htmlRender *HtmlTemplateRender) Render(write http.ResponseWriter, result Result) {
	write.Header().Set("Content-Type", "text/html;charset=utf8")

	write.WriteHeader(result.code)
	// 使用模板渲染
	if err := htmlRender.getHtmlTemplates().ExecuteTemplate(write, result.htmlTplName, result.data); err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		write.Write([]byte(err.Error()))
	}
}

func (htmlRender *HtmlTemplateRender) getHtmlTemplates() *template.Template {
	return htmlRender.HtmlTemplates
}

func (htmlRender *HtmlTemplateRender) GetFuncMap() template.FuncMap {
	return htmlRender.FuncMap
}
func (htmlRender *HtmlTemplateRender) SetHtmlTemplates(htmlTemplates *template.Template) {
	htmlRender.HtmlTemplates = htmlTemplates
}
func (htmlRender *HtmlTemplateRender) SetFuncMap(funcMap template.FuncMap) {
	htmlRender.FuncMap = funcMap
}
