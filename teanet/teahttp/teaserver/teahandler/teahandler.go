package teahandler

import (
	"context"
	"net/http"

	"github.com/tea-frame-go/tea/teaerrors"
	"github.com/tea-frame-go/tea/teanet/teahttp/teaserver"
	"github.com/tea-frame-go/tea/teatypes"
	"github.com/tea-frame-go/tea/teavalid"
)

func Handler[Req teavalid.Validator[T], T any, Res any](
	logicHandler func(ctx context.Context, req Req) (res teatypes.Result[Res], err teaerrors.Error),
	errorHandler func(r *teaserver.Request, err teaerrors.Error),
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
		res, err := logicHandler(r.Context(), req)
		switch {
		case err != nil:
			errorHandler(r, err)
			r.Response.WriteStatus(http.StatusInternalServerError)
			return
		case res.Error != nil:
			r.Response.WriteJson(responseError{res.Error})
			return
		default:
			r.Response.WriteJson(response{res.Value})
			return
		}
	}
}

type responseError struct {
	Error error
}

type response struct {
	Data any
}
