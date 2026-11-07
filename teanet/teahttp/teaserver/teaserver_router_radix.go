package teaserver

import (
	"fmt"
	"slices"
	"strings"
)

type routerRadixNode struct {
	path     string
	indexes  []byte
	children []*routerRadixNode
	handler  HandlerFunc
}

func newRouterRadixNode() routerRadixNode {
	return routerRadixNode{
		path:     "/",
		indexes:  nil,
		children: nil,
		handler:  nil,
	}
}

func (n *routerRadixNode) insert(path string, handler HandlerFunc) {
	originPath := path
	for {
		length := n.commonPrefixLength(path)
		if length == len(n.path) {
			if length == len(path) {
				if n.handler != nil {
					panic(fmt.Errorf("teahttp: router path `%s` conflict", originPath))
				}
				n.handler = handler
				return
			}
			path = strings.Clone(path[length:])
			if child := n.child(path[0]); child != nil {
				n = child
				continue
			}
			newNode := &routerRadixNode{
				path:     path,
				indexes:  nil,
				children: nil,
				handler:  handler,
			}
			n.indexes = append(n.indexes, newNode.path[0])
			n.children = append(n.children, newNode)
			return
		}
		node := &routerRadixNode{
			path:     n.path[length:],
			indexes:  n.indexes,
			children: n.children,
			handler:  n.handler,
		}
		n.path = n.path[:length]
		n.indexes = []byte{node.path[0]}
		n.children = []*routerRadixNode{node}
		n.handler = nil
		if length == len(path) {
			n.handler = handler
			return
		}
		path = strings.Clone(path[length:])
		newNode := &routerRadixNode{
			path:     path,
			indexes:  nil,
			children: nil,
			handler:  handler,
		}
		n.indexes = append(n.indexes, newNode.path[0])
		n.children = append(n.children, newNode)
		return
	}
}

func (n *routerRadixNode) search(path string) HandlerFunc {
	for {
		length := n.commonPrefixLength(path)
		if length == len(n.path) {
			if length == len(path) {
				return n.handler
			}
			path = path[length:]
			if child := n.child(path[0]); child != nil {
				n = child
				continue
			}
			return nil
		}
		return nil
	}
}

func (n *routerRadixNode) clip() {
	n.indexes = slices.Clone(n.indexes)
	n.children = slices.Clone(n.children)
	for _, child := range n.children {
		child.clip()
	}
}

func (n *routerRadixNode) commonPrefixLength(path string) int {
	length := 1
	lengthMax := min(len(n.path), len(path))
	for length < lengthMax && n.path[length] == path[length] {
		length++
	}
	return length
}

func (n *routerRadixNode) child(head byte) *routerRadixNode {
	for i, index := range n.indexes {
		if index == head {
			return n.children[i]
		}
	}
	return nil
}
