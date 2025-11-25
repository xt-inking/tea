package teamiddleware

import (
	"net/http"
	"slices"

	"github.com/tea-frame-go/tea/teanet/teahttp/teaserver"
)

func Cors(allowOrigins []string, allowHeaders string) func(next teaserver.HandlerFunc) teaserver.HandlerFunc {
	return func(next teaserver.HandlerFunc) teaserver.HandlerFunc {
		return func(r *teaserver.Request) {
			origin := r.Raw.Header.Get("Origin")
			if origin == "" {
				next(r)
				return
			}
			if !slices.Contains(allowOrigins, origin) {
				return
			}
			r.Response.Header().Set("Access-Control-Allow-Origin", origin)
			r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Raw.Method == http.MethodOptions {
				r.Response.Header().Set("Access-Control-Max-Age", "7200")
				r.Response.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
				r.Response.Header().Set("Access-Control-Allow-Headers", allowHeaders)
				return
			}
			next(r)
		}
	}
}
