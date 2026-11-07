package teaserver

import (
	"slices"
)

type router struct {
	radix routerRadixNode
}

func newRouter() router {
	return router{
		radix: newRouterRadixNode(),
	}
}

func (r *router) insert(path string, handler HandlerFunc) {
	r.radix.insert(path, handler)
}

func (r *router) search(path string) HandlerFunc {
	return r.radix.search(path)
}

func (r *router) Tree(path string, f func(t *RouterTree)) {
	f(newRouterTree(cleanPath(path), nil, r))
	r.radix.clip()
}

type RouterTree struct {
	path       string
	middleware []func(next HandlerFunc) HandlerFunc
	router     *router
}

func newRouterTree(path string, middleware []func(next HandlerFunc) HandlerFunc, router *router) *RouterTree {
	return &RouterTree{
		path:       path,
		middleware: middleware,
		router:     router,
	}
}

func (t *RouterTree) Tree(path string, f func(t *RouterTree)) {
	f(newRouterTree(t.path+cleanPath(path), slices.Clone(t.middleware), t.router))
}

func (t *RouterTree) Middleware(middleware ...func(next HandlerFunc) HandlerFunc) {
	t.middleware = append(t.middleware, middleware...)
}

func (t *RouterTree) Post(path string, handler HandlerFunc) {
	t.register(path, handler)
}

func (t *RouterTree) register(path string, handler HandlerFunc) {
	path = t.path + cleanPath(path)
	if path == "" {
		path = "/"
	}
	for _, m := range slices.Backward(t.middleware) {
		handler = m(handler)
	}
	t.router.insert(path, handler)
}

func cleanPath(path string) string {
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
