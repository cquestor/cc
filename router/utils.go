package router

import (
	"fmt"
	"strings"
)

// GetCNodeType 判断路由节点的类型
func GetCNodeType(v string) TypeNode {
	if isNamed(v) {
		return NamedNode
	}
	if isWildcard(v) {
		return WildcardNode
	}
	return AbsoluteNode
}

// isNamed 判断是否是命名路由
func isNamed(v string) bool {
	return strings.HasPrefix(v, ":")
}

// isWildcard 判断是否是通配符路由
func isWildcard(v string) bool {
	return strings.HasPrefix(v, "*")
}

// ParsePattern 获取路径节点
func ParsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
			if strings.HasPrefix(part, "*") {
				break
			}
		}
	}
	return result
}

// CheckValid 检查路由冲突
func CheckValid(parent *CNode, child *CNode) {
	for _, each := range parent.Children {
		if each == child || each.Type == AbsoluteNode {
			continue
		}
		if each.Type == child.Type && each.Pattern != "" {
			panic(fmt.Sprintf("route conflict(%s): %s", each.Pattern, child.Pattern))
		}
	}
}
