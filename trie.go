package goee

import "strings"

// 前缀树实现动态路由匹配
// 定义数结构
type node struct {
	pattern  string  // 待匹配的路由,例如 /p/:lang
	part     string  // 路由中的一部分,例如:lang
	children []*node // 子节点，例如[doc,info]
	isWild   bool    // 是否精确匹配，part含有：或*时为true
	isStatic bool    // 是否静态文件路由
}

// 匹配子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配所有子节点
func (n *node) matchChildren(part string) []*node {
	// 定义一个节点slice
	nodes := make([]*node, 0)

	for _, child := range n.children {
		// 如果匹配上，加入切片
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入节点
func (n *node) insert(pattern string, parts []string, isStatic bool, height int) {
	// 如果树的高度和节点路由部分长度一致，说明已到叶子节点，直接返回
	if len(parts) == height {
		// 如果该节点pattern为空，则赋值，否则说明已存在，不允许覆盖
		if n.pattern == "" {
			n.pattern = pattern
		}
		return
	}

	// 查找当前高度和对应部分路由，如果不存在节点则插入节点，并继续子节点的构造
	part := parts[height]
	// 查找时候存在该部分的子节点
	child := n.matchChild(part)
	// 不存在则创建插入
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*', isStatic: isStatic}
		n.children = append(n.children, child)
	}
	// 继续插入子节点
	child.insert(pattern, parts, isStatic, height+1)
}

// 查找节点
func (n *node) search(parts []string, height int) *node {
	// 如果已查找的叶子节点或者通配符节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果节点的pattern为空，说明该节点异常，否则直接返回
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 查找每部分
	part := parts[height]
	children := n.matchChildren(part)

	// 遍历查找每一个符合的子节点
	for _, child := range children {
		result := child.search(parts, height+1)
		// 查找到之后直接返回
		if result != nil {
			return result
		}
	}

	// 没找到
	return nil
}
