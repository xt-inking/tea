package teamiddleware

import (
	"net/http"

	"github.com/tea-frame-go/tea/teaerrors"
	"github.com/tea-frame-go/tea/tealog"
	"github.com/tea-frame-go/tea/teanet/teahttp/teaserver"
)

func Recover(s *teaserver.Server) func(next teaserver.HandlerFunc) teaserver.HandlerFunc {
	logger := tealog.New(
		tealog.NewRecordHandlerText(),
		tealog.NewWriterCloserFile(tealog.FileDir, "http-recover"),
	)
	s.Loggers(logger)
	return func(next teaserver.HandlerFunc) teaserver.HandlerFunc {
		return func(r *teaserver.Request) {
			defer func() {
				if err := recover(); err != nil {
					e := teaerrors.NewAny(err, 2)
					logger.Error(r.Context(), e.ErrorStack())
					r.Response.WriteStatus(http.StatusInternalServerError)
				}
			}()
			next(r)
		}
	}
}
