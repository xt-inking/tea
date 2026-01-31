package teaserver

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/tea-frame-go/tea/teaencoding/teajson"
)

type Request struct {
	Raw      *http.Request
	Response *response
	server   *Server
	ctx      context.Context
}

func (request *Request) Context() context.Context {
	return request.ctx
}

func (request *Request) SetContext(ctx context.Context) {
	request.ctx = ctx
}

func (request *Request) Form() (*multipart.Form, error) {
	const defaultMaxMemory = 32 << 20 // 32 MB
	err := request.Raw.ParseMultipartForm(defaultMaxMemory)
	return request.Raw.MultipartForm, err
}

func (request *Request) Json(v any) error {
	err := teajson.Decode(request.Raw.Body, v)
	return err
}
