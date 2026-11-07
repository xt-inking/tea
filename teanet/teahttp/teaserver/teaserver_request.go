package teaserver

import (
	"context"
	"encoding/json/v2"
	"mime/multipart"
	"net/http"

	"github.com/xt-inking/tea/internal/profile"
	"github.com/xt-inking/tea/teaerrors"
)

type Request struct {
	Raw      *http.Request
	Response response
	server   *Server
	ctx      context.Context
}

func (r *Request) Form() (*multipart.Form, error) {
	const defaultMaxMemory = 32 << 20
	err := r.Raw.ParseMultipartForm(defaultMaxMemory)
	if profile.Debug {
		if err != nil {
			r.server.logger.Error(r.Ctx(), teaerrors.New(err).Error())
		}
	}
	return r.Raw.MultipartForm, err
}

func (r *Request) Json(out any) error {
	err := json.UnmarshalRead(r.Raw.Body, out)
	if profile.Debug {
		if err != nil {
			r.server.logger.Error(r.Ctx(), teaerrors.New(err).Error())
		}
	}
	return err
}

func (r *Request) Ctx() context.Context {
	return r.ctx
}

func (r *Request) SetCtx(ctx context.Context) {
	r.ctx = ctx
}
