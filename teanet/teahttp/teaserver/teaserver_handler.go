package teaserver

import (
	"net/http"
	"sync"

	"github.com/tea-frame-go/tea/teaerrors"
)

type HandlerFunc func(r *Request)

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	request := requestGet(r, w, s)
	defer func() {
		if err := recover(); err != nil {
			e := teaerrors.NewAny(err, 2)
			s.logger.Error(request.Context(), e.ErrorStack())
			request.Response.WriteStatus(http.StatusInternalServerError)
		}
		requestPut(request)
	}()
	if handler := s.router.search(r.URL.Path); handler != nil {
		handler(request)
	} else {
		request.Response.WriteStatus(http.StatusNotFound)
	}
}

var requestPool = sync.Pool{
	New: func() any {
		request := &Request{
			Raw: nil,
			Response: &response{
				writer:  nil,
				Request: nil,
			},
			server: nil,
			ctx:    nil,
		}
		request.Response.Request = request
		return request
	},
}

func requestGet(r *http.Request, w http.ResponseWriter, s *server) *Request {
	request := requestPool.Get().(*Request)
	request.Raw = r
	request.Response.writer = w
	request.server = s
	request.ctx = r.Context()
	return request
}

func requestPut(request *Request) {
	request.Raw = nil
	request.Response.writer = nil
	request.server = nil
	request.ctx = nil
	requestPool.Put(request)
}
