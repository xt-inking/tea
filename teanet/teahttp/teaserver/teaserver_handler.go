package teaserver

import (
	"net/http"
	"sync"
)

type HandlerFunc func(r *Request)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := requestPoolGet(r, w, s)
	if handler := s.router.search(request.Raw.URL.Path); handler != nil {
		handler(request)
	} else {
		request.Response.WriteStatus(http.StatusNotFound)
	}
	requestPoolPut(request)
}

var requestPool = sync.Pool{
	New: func() any {
		request := &Request{
			Raw: nil,
			Response: response{
				Writer:  nil,
				request: nil,
			},
			server: nil,
			ctx:    nil,
		}
		request.Response.request = request
		return request
	},
}

func requestPoolGet(r *http.Request, w http.ResponseWriter, s *Server) *Request {
	request := requestPool.Get().(*Request)
	request.Raw = r
	request.Response.Writer = w
	request.server = s
	request.ctx = r.Context()
	return request
}

func requestPoolPut(request *Request) {
	request.Raw = nil
	request.Response.Writer = nil
	request.server = nil
	request.ctx = nil
	requestPool.Put(request)
}
