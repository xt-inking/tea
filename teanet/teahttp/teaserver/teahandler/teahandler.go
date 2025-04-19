package teahandler

import (
	"context"
	"net/http"

	"github.com/tea-frame-go/tea/teaerrors"
	"github.com/tea-frame-go/tea/teaerrors/teaerrorcode"
	"github.com/tea-frame-go/tea/teanet/teahttp/teaserver"
	"github.com/tea-frame-go/tea/teavalid"
)

func Handler[Req teavalid.Validator[T], T any, Res any](
	f func(ctx context.Context, req Req) (res Res, err error),
	errorHandler func(r *teaserver.Request, e teaerrors.ErrorStack),
) teaserver.HandlerFunc {
	return func(r *teaserver.Request) {
		req := new(T)
		if err := r.Json(req); err != nil {
			r.Response.WriteStatusText(http.StatusBadRequest, "json decode error")
			return
		}
		if err := Req(req).Validate(); err != nil {
			r.Response.WriteStatusText(http.StatusBadRequest, err.Error())
			return
		}
		res, err := f(r.Context(), req)
		if err != nil {
			if e, ok := err.(teaerrorcode.ErrorCode); ok {
				r.Response.WriteJson(e)
				return
			}
			e, ok := err.(teaerrors.ErrorStack)
			if !ok {
				e = teaerrors.New(err, 0)
			}
			errorHandler(r, e)
			r.Response.WriteStatus(http.StatusInternalServerError)
			return
		}
		r.Response.WriteJson(Response(res))
	}
}

func Response[T any](res T) response[T] {
	return response[T]{
		Code: 0,
		Data: res,
	}
}

type response[T any] struct {
	Code uint8
	Data T
}
