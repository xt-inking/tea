package teaserver

import (
	"net/http"
	"unsafe"

	"github.com/tea-frame-go/tea/teaencoding/teajson"
	"github.com/tea-frame-go/tea/teaerrors"
)

type response struct {
	Writer  http.ResponseWriter
	Request *Request
}

func (response *response) Header() http.Header {
	return response.Writer.Header()
}

func (response *response) WriteStatus(code int) {
	text := http.StatusText(code)
	response.WriteStatusText(code, text)
}

func (response *response) WriteStatusText(code int, text string) {
	response.Writer.WriteHeader(code)
	response.Write(unsafe.Slice(unsafe.StringData(text), len(text)))
}

func (response *response) Write(data []byte) {
	_, err := response.Writer.Write(data)
	if err != nil {
		e := teaerrors.New(err, 0)
		response.Request.server.logger.Error(response.Request.Context(), e.ErrorStack())
	}
}

func (response *response) WriteJson(data any) {
	response.Header().Set("Content-Type", "application/json")
	err := teajson.Encode(response.Writer, data)
	if err != nil {
		e := teaerrors.New(err, 0)
		response.Request.server.logger.Error(response.Request.Context(), e.ErrorStack())
	}
}
