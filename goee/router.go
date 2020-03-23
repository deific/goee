package goee

import (
	"net/http"
	"strings"
)

// 路由结构体
type router struct {
	roots    map[string]*node // 路由树,例如：roots["GET"],roots["POST"]
	handlers map[string]HandlerFunc
}

// 创建一个路由器
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析请求路径
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 如果是通配符,不用解析直接跳出
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	// 解析
	parts := parsePattern(pattern)

	key := method + ":" + pattern
	// 获取路由根节点,不存在则初始化一个
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 插入路由
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 根据路由获取路由
func (r *router) getRouter(method string, pattern string) (*node, map[string]string) {
	// 解析
	searchParts := parsePattern(pattern)
	// 声明用于保存url中的参数
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 查找路由节点
	n := root.search(searchParts, 0)
	if n != nil {
		// 解析url参数
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 如果配到的节点路由部分以:开头，说明是参数
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// 如果匹配到的部分是*符号，说明后续部分都是参数值
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[1:], "/")
				break
			}
		}
		return n, params
	}
	// 没找到
	return nil, nil
}

// 处理
func (r *router) handle(c *Context) {

	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + ":" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	}
}
