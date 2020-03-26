package render

import (
	"html/template"
)

type Enum int16

type Render interface {
	GetHtmlTemplates() *template.Template
	GetFuncMap() template.FuncMap
	SetHtmlTemplates(template *template.Template)
	SetFuncMap(funcMap template.FuncMap)
}

// 渲染器管理器
type GoeeRenderManager struct {
	renderType Enum // 渲染类型
	// html模板
	HtmlTemplates *template.Template
	FuncMap       template.FuncMap
}

func New() Render {
	return &GoeeRenderManager{}
}

func (renderManager *GoeeRenderManager) GetHtmlTemplates() *template.Template {
	return renderManager.HtmlTemplates
}

func (renderManager *GoeeRenderManager) GetFuncMap() template.FuncMap {
	return renderManager.FuncMap
}
func (renderManager *GoeeRenderManager) SetHtmlTemplates(htmlTemplates *template.Template) {
	renderManager.HtmlTemplates = htmlTemplates
}
func (renderManager *GoeeRenderManager) SetFuncMap(funcMap template.FuncMap) {
	renderManager.FuncMap = funcMap
}
