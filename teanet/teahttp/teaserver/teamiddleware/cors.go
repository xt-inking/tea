package teamiddleware

import (
	"net/http"
	"slices"

	"github.com/tea-frame-go/tea/teaconfig"
	"github.com/tea-frame-go/tea/teanet/teahttp/teaserver"
)

func Cors(config *teaconfig.CorsConfig) func(next teaserver.HandlerFunc) teaserver.HandlerFunc {
	return func(next teaserver.HandlerFunc) teaserver.HandlerFunc {
		return func(r *teaserver.Request) {
			origin := r.Raw.Header.Get("Origin")
			if origin == "" {
				next(r)
				return
			}
			if !slices.Contains(config.AllowOrigins, origin) {
				return
			}
			r.Response.Header().Set("Access-Control-Allow-Origin", origin)
			r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Raw.Method == http.MethodOptions {
				r.Response.Header().Set("Access-Control-Max-Age", "86400")
				r.Response.Header().Set("Access-Control-Allow-Methods", "POST")
				r.Response.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders)
				return
			}
			next(r)
		}
	}
}
