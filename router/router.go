package router

import (
	"sort"
)

// TypeNode 节点类型
type TypeNode int

// IRouter 路由接口
type IRouter interface {
	AddRoute(method, pattern string)                             // 添加路由
	GetRoute(method, pattern string) (string, map[string]string) // 获取路由和路由参数
}

// CRouter 路由
type CRouter struct {
	Roots map[string]*CNode
}

// CNode 路由节点
type CNode struct {
	Pattern  string
	Part     string
	Type     TypeNode
	Children []*CNode
}

// CNodes 节点组，用于排序
type CNodes []*CNode

const (
	WildNode     TypeNode = iota // 通配符路由
	DynamicNode                  // 动态路由
	AbsoluteNode                 // 绝对路由
)

// NewRouter 构造路由
func NewRouter() *CRouter {
	return &CRouter{
		Roots: make(map[string]*CNode),
	}
}

// NewNode 构造路由节点
func NewNode(part string) *CNode {
	return &CNode{
		Part:     part,
		Type:     CheckNodeType(part),
		Children: make([]*CNode, 0),
	}
}

// AddRoute 添加路由
func (router *CRouter) AddRoute(method, pattern string) {
	if _, ok := router.Roots[method]; !ok {
		router.Roots[method] = NewNode("/")
	}
	router.Roots[method].insert(pattern, ParsePattern(pattern), router.Roots[method])
}

// GetRoute 获取路由节点及路由参数
func (router *CRouter) GetRoute(method string, pattern string) (string, map[string]string) {
	searchParts := ParsePattern(pattern)
	root, ok := router.Roots[method]
	if !ok {
		return "", nil
	}
	if n := root.search(searchParts); n != nil {
		parts := ParsePattern(n.Pattern)
		return n.Pattern, ParseParams(parts, searchParts)
	}
	return "", nil
}

// insert 路由节点插入
func (n *CNode) insert(pattern string, parts []string, parent *CNode) {
	if len(parts) == 0 {
		n.Pattern = pattern
		CheckValid(parent, n)
		return
	}
	part := parts[0]
	child := n.matchChild(part)
	if child == nil {
		child = NewNode(part)
		n.Children = append(n.Children, child)
		sort.Sort(sort.Reverse(CNodes(n.Children)))
	}
	child.insert(pattern, parts[1:], n)
}

// search 路由节点查找
func (n *CNode) search(parts []string) *CNode {
	if len(parts) == 0 || isWild(n.Part) {
		if n.Pattern == "" {
			return nil
		}
		return n
	}
	part := parts[0]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts[1:])
		if result != nil {
			return result
		}
	}
	return nil
}

// matchChild 匹配路由节点，用于插入
func (n *CNode) matchChild(part string) *CNode {
	for _, child := range n.Children {
		if child.Part == part {
			return child
		}
	}
	return nil
}

// matchChildren 匹配路由节点，用于查找
func (n *CNode) matchChildren(part string) []*CNode {
	nodes := make([]*CNode, 0)
	for _, child := range n.Children {
		if child.Part == part || child.Type < 2 {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (nodes CNodes) Len() int {
	return len(nodes)
}

func (nodes CNodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func (nodes CNodes) Less(i, j int) bool {
	return nodes[i].Type < nodes[j].Type
}
