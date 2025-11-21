package teaclient

import (
	"net/http"

	"github.com/tea-frame-go/tea/teaencoding/teajson"
)

type Response struct {
	Raw *http.Response
}

func newResponse(raw *http.Response) *Response {
	resp := &Response{
		Raw: raw,
	}
	return resp
}

func (r *Response) Close() error {
	return r.Raw.Body.Close()
}

func (r *Response) Json(v any) error {
	return teajson.Decode(r.Raw.Body, v)
}
