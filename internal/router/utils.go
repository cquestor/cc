package router

import (
	"fmt"
	"strings"
)

// CheckNodeType 判断路由节点类型
func CheckNodeType(part string) TypeNode {
	if isWild(part) {
		return WildNode
	}
	if isDynamic(part) {
		return DynamicNode
	}
	return AbsoluteNode
}

// ParasePattern 解析路由
func ParsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
			if isWild(part) {
				break
			}
		}
	}
	return result
}

// ParseParams 解析路由参数
func ParseParams(parts []string, searchParts []string) map[string]string {
	params := make(map[string]string)
	for index, part := range parts {
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = searchParts[index]
		}
		if strings.HasPrefix(part, "*") && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	return params
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

// isWild 判断是否为通配符路由节点
func isWild(part string) bool {
	if len(part) <= 1 {
		return false
	}
	return strings.HasPrefix(part, "*")
}

// isDynamic 判断是否为动态路由节点
func isDynamic(part string) bool {
	if len(part) <= 1 {
		return false
	}
	return strings.HasPrefix(part, ":")
}
