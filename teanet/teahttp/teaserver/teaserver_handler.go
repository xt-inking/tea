package teaserver

import (
	"net/http"
	"sync"
)

type HandlerFunc func(r *Request)

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	request := requestGet(r, w, s)
	if handler := s.router.search(r.URL.Path); handler != nil {
		handler(request)
	} else {
		request.Response.WriteStatus(http.StatusNotFound)
	}
	requestPut(request)
}

var requestPool = sync.Pool{
	New: func() any {
		request := &Request{
			Raw: nil,
			Response: &response{
				Writer:  nil,
				Request: nil,
			},
			server: nil,
			ctx:    nil,
		}
		request.Response.Request = request
		return request
	},
}

func requestGet(r *http.Request, w http.ResponseWriter, s *Server) *Request {
	request := requestPool.Get().(*Request)
	request.Raw = r
	request.Response.Writer = w
	request.server = s
	request.ctx = r.Context()
	return request
}

func requestPut(request *Request) {
	request.Raw = nil
	request.Response.Writer = nil
	request.server = nil
	request.ctx = nil
	requestPool.Put(request)
}
