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
	server   *server
	ctx      context.Context
	m        map[string]any
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

func (request *Request) Json() error {
	err := teajson.Decode(request.Raw.Body, &request.m)
	return err
}

func (request *Request) Map() map[string]any {
	return request.m
}
