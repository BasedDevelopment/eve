package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Logger uses zerolog to log information about each request (log level = INFO)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		reqId := middleware.GetReqID(r.Context())

		defer func() {
			log.Info().
				Str("reqId", reqId).
				Str("method", r.Method).
				Str("host", r.Host).
				Str("client", r.RemoteAddr).
				Str("page", r.RequestURI).
				Str("protocol", r.Proto).
				Str("user-agent", r.UserAgent()).
				Int64("duration(μs)", time.Since(t).Microseconds()).
				Int("status", ww.Status()).
				Str("bytes_in", r.Header.Get("Content-Length")).
				Int("bytes_out", ww.BytesWritten()).
				Msg("HTTP Request")
		}()

		next.ServeHTTP(ww, r)
	})
}
