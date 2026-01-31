package teaserver

import (
	"fmt"
)

type routerRadixNode struct {
	path     string
	handler  HandlerFunc
	indexes  string
	children []*routerRadixNode
}

func newRouterRadixNode() *routerRadixNode {
	node := &routerRadixNode{
		path:     "/",
		handler:  nil,
		indexes:  "",
		children: nil,
	}
	return node
}

func (node *routerRadixNode) insert(path string, handler HandlerFunc) {
	rawPath := path
	for {
		length := node.commonPrefixLength(path)
		if length < len(node.path) {
			newNode := &routerRadixNode{
				path:     node.path[length:],
				handler:  node.handler,
				indexes:  node.indexes,
				children: node.children,
			}
			node.path = node.path[:length]
			node.handler = nil
			node.indexes = newNode.path[0:1]
			node.children = []*routerRadixNode{newNode}
		}
		if length < len(path) {
			if childNode := node.childNode(path[length]); childNode != nil {
				node = childNode
				path = path[length:]
				continue
			}
			newNode := &routerRadixNode{
				path:     path[length:],
				handler:  handler,
				indexes:  "",
				children: nil,
			}
			node.indexes += newNode.path[0:1]
			node.children = append(node.children, newNode)
			return
		}
		if node.handler != nil {
			panic(fmt.Sprintf("path `%s` conflicts", rawPath))
		}
		node.handler = handler
		return
	}
}

func (node *routerRadixNode) search(path string) HandlerFunc {
	for {
		length := node.commonPrefixLength(path)
		if length == len(node.path) {
			if length == len(path) {
				return node.handler
			}
			if childNode := node.childNode(path[length]); childNode != nil {
				node = childNode
				path = path[length:]
				continue
			}
			return nil
		}
		return nil
	}
}

func (node *routerRadixNode) commonPrefixLength(path string) int {
	length := 1
	max := min(len(node.path), len(path))
	for length < max && node.path[length] == path[length] {
		length++
	}
	return length
}

func (node *routerRadixNode) childNode(head byte) *routerRadixNode {
	for i := range len(node.indexes) {
		if node.indexes[i] == head {
			return node.children[i]
		}
	}
	return nil
}
