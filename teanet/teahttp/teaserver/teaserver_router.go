package teaserver

import (
	"slices"
)

type router struct {
	radix *routerRadixNode
}

func newRouter() *router {
	r := &router{
		radix: newRouterRadixNode(),
	}
	return r
}

func (r *router) Tree(path string, f func(t *RouterTree)) {
	f(newRouterTree(realPath(path), nil, r))
}

func (r *router) insert(path string, handler HandlerFunc) {
	r.radix.insert(path, handler)
}

func (r *router) search(path string) HandlerFunc {
	return r.radix.search(path)
}

type RouterTree struct {
	path       string
	middleware []func(next HandlerFunc) HandlerFunc
	router     *router
}

func newRouterTree(path string, middleware []func(next HandlerFunc) HandlerFunc, r *router) *RouterTree {
	t := &RouterTree{
		path:       path,
		middleware: middleware,
		router:     r,
	}
	return t
}

func (t *RouterTree) Tree(path string, f func(t *RouterTree)) {
	f(newRouterTree(t.path+realPath(path), slices.Clone(t.middleware), t.router))
}

func (t *RouterTree) Middleware(middleware ...func(next HandlerFunc) HandlerFunc) {
	t.middleware = append(t.middleware, middleware...)
}

func (t *RouterTree) Post(path string, handler HandlerFunc) {
	t.register(path, handler)
}

func (t *RouterTree) register(path string, handler HandlerFunc) {
	path = t.path + realPath(path)
	if path == "" {
		path = "/"
	}
	for i := len(t.middleware) - 1; i >= 0; i-- {
		handler = t.middleware[i](handler)
	}
	t.router.insert(path, handler)
}

func realPath(path string) string {
	switch {
	case path == "":
		return ""
	case path == "/":
		return ""
	case path[0] == '/':
		return path
	default:
		return "/" + path
	}
}
