package router

import (
	"sort"
	"strings"
)

// CRouter 路由
type CRouter struct {
	Roots map[string]*CNode
}

// CNode 路由节点
type CNode struct {
	Current  string
	Type     TypeNode
	Children []*CNode
	Pattern  string
}

// CNodes 节点排序
type CNodes []*CNode

// TypeNode 节点类型
type TypeNode int

const (
	WildcardNode TypeNode = iota // 通配符路由
	NamedNode                    // 命名路由
	AbsoluteNode                 // 绝对路由
)

// NewCRouter 新建路由
func NewCRouter() *CRouter {
	return &CRouter{
		Roots: make(map[string]*CNode),
	}
}

// NewCNode 新建路由节点
func NewCNode(current string, typ TypeNode) *CNode {
	return &CNode{
		Current:  current,
		Type:     typ,
		Children: make([]*CNode, 0),
	}
}

// AddRoute 添加路由
func (r *CRouter) AddRoute(method, pattern string) {
	_, ok := r.Roots[method]
	if !ok {
		r.Roots[method] = NewCNode("/", AbsoluteNode)
	}
	r.Roots[method].insert(pattern, ParsePattern(pattern), r.Roots[method])
}

// GetRoute 获取路由
func (r *CRouter) GetRoute(method string, pattern string) (*CNode, map[string]string) {
	searchParts := ParsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.Roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts)
	if n != nil {
		parts := ParsePattern(n.Pattern)
		for index, part := range parts {
			if strings.HasPrefix(part, ":") {
				params[part[1:]] = searchParts[index]
			}
			if strings.HasPrefix(part, "*") && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// insert 插入节点
func (n *CNode) insert(pattern string, parts []string, parent *CNode) {
	if len(parts) == 0 {
		n.Pattern = pattern
		CheckValid(parent, n)
		return
	}
	part := parts[0]
	child := n.matchChild(part)
	if child == nil {
		child = NewCNode(part, GetCNodeType(part))
		n.Children = append(n.Children, child)
		sort.Sort(sort.Reverse(CNodes(n.Children)))
	}
	child.insert(pattern, parts[1:], n)
}

// matchChild 匹配节点
func (n *CNode) matchChild(part string) *CNode {
	for _, child := range n.Children {
		if child.Current == part {
			return child
		}
	}
	return nil
}

// search 查找匹配路由
func (n *CNode) search(parts []string) *CNode {
	if len(parts) == 0 || strings.HasPrefix(n.Current, "*") {
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

// matchChildren 匹配节点
func (n *CNode) matchChildren(part string) []*CNode {
	nodes := make([]*CNode, 0)
	for _, child := range n.Children {
		if child.Current == part || child.Type < 2 {
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
