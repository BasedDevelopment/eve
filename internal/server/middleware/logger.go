/*
 * eve - management toolkit for libvirt servers
 * Copyright (C) 2022-2023  BNS Services LLC

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Logger uses zerolog to log information about each request (log level = INFO)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		reqId := middleware.GetReqID(r.Context())

		// Copy the request body so we can log it
		var (
			bodyStr   string
			bodyBytes []byte
			err       error
		)

		// Log r.body if the request does not contain a password
		if !(r.Body == nil || strings.Contains(r.RequestURI, "login") || strings.Contains(r.RequestURI, "/admin/user")) {
			bodyBytes, err = io.ReadAll(r.Body)
			bodyStr = string(bodyBytes)

			if err != nil {
				util.WriteError(w, r, err, http.StatusInternalServerError, "failed to read request body")
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		defer func() {
			logger := log.With().
				Str("reqId", reqId).
				Str("method", r.Method).
				Str("host", r.Host).
				Str("client", r.RemoteAddr).
				Str("page", r.RequestURI).
				Str("protocol", r.Proto).
				Str("user-agent", r.UserAgent()).
				Float64("duration(ms)", float64(time.Since(t).Nanoseconds())/1000000.0).
				Int("status", ww.Status()).
				Int("bytes_in", len(bodyBytes)).
				Int("bytes_out", ww.BytesWritten()).
				Logger()

			if ww.Status() >= 500 {
				logger.Error().Str("body", bodyStr).Msg("Request failed")
			} else {
				logger.Info().Msg("Request completed")
			}
		}()

		next.ServeHTTP(ww, r)
	})
}
