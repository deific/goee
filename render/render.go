package render

import (
	"html/template"
	"log"
	"net/http"
)

type RenderType string

const (
	STRING       RenderType = "string"
	JSON         RenderType = "json"
	HTML         RenderType = "html"
	HTMLTEMPLATE RenderType = "htmltemplate"
)

type Render interface {
	Render(write http.ResponseWriter, result Result)
}

// 结果类型
type Result struct {
	code        int
	data        interface{} // 响应码
	htmlTplName string
}

// 渲染器管理器
type GoeeRenderManager struct {
	renders map[string]Render // 渲染类型
}

// 初始化
func New() *GoeeRenderManager {
	log.Println("初始化渲染管理器...")
	return &GoeeRenderManager{renders: map[string]Render{
		string(HTML):         &HtmlRender{},
		string(HTMLTEMPLATE): &HtmlTemplateRender{},
	}}
}

// 渲染
func (manager *GoeeRenderManager) DoRender(renderType RenderType, write http.ResponseWriter, code int, data interface{}, htmlTemplateName string) {
	result := Result{code: code, data: data, htmlTplName: htmlTemplateName}
	render := manager.renders[string(renderType)]
	render.Render(write, result)
}

func (manager *GoeeRenderManager) SetHtmlTemplates(htmlTemplates *template.Template) {
	render := manager.renders[string(HTMLTEMPLATE)]
	// 类型断言，左侧只能是interface类型才可以，因此需要使用interface包装一下
	htmlTplRender, isTemplateRender := interface{}(render).(HtmlTemplateRender)
	if !isTemplateRender {
		htmlTplRender.HtmlTemplates = htmlTemplates
	}
}
func (manager *GoeeRenderManager) SetFuncMap(funcMap template.FuncMap) {
	render := manager.renders[string(HTMLTEMPLATE)]
	htmlTplRender, isTemplateRender := interface{}(render).(HtmlTemplateRender)
	if isTemplateRender {
		htmlTplRender.FuncMap = funcMap
	}
}

func (manager *GoeeRenderManager) GetFuncMap() template.FuncMap {
	render := manager.renders[string(HTMLTEMPLATE)]
	htmlTplRender, isTemplateRender := interface{}(render).(HtmlTemplateRender)
	if !isTemplateRender {
		return htmlTplRender.FuncMap
	}
	return nil
}
