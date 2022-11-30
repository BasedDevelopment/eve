package middlewares

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.Info().
				// TODO: request ids?
				Str("method", r.Method).
				Str("host", r.Host).
				Str("client", r.RemoteAddr).
				Str("page", r.RequestURI).
				Str("protocol", r.Proto).
				Msg("HTTP Request")
		}()

		next.ServeHTTP(w, r)
	})
}
