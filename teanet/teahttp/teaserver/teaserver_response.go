package teaserver

import (
	"encoding/json/v2"
	"net/http"
	"unsafe"

	"github.com/xt-inking/tea/teaerrors"
)

type response struct {
	Writer  http.ResponseWriter
	request *Request
}

func (r *response) Header() http.Header {
	return r.Writer.Header()
}

func (r *response) WriteStatus(code int) {
	r.WriteStatusText(code, http.StatusText(code))
}

func (r *response) WriteStatusText(code int, text string) {
	r.Writer.WriteHeader(code)
	r.Write(unsafe.Slice(unsafe.StringData(text), len(text)))
}

func (r *response) Write(data []byte) {
	_, err := r.Writer.Write(data)
	if err != nil {
		r.request.server.logger.Error(r.request.Ctx(), teaerrors.New(err).Error())
	}
}

func (r *response) WriteJson(data any) {
	r.Header().Set("Content-Type", "application/json")
	err := json.MarshalWrite(r.Writer, data)
	if err != nil {
		r.request.server.logger.Error(r.request.Ctx(), teaerrors.New(err).Error())
	}
}
