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

package util

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func WriteResponse[R any](r R, w http.ResponseWriter, status int) error {
	json, err := json.Marshal(r)

	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return nil
}

// Writes error response back to the client
// Logs the error if it is an actual server error (err != nil)
func WriteError(w http.ResponseWriter, r *http.Request, e error, s int, m string) {
	if e != nil && s == http.StatusInternalServerError {
		log.Error().
			Err(e).
			Str("message", m).
			Msg("Request Error")
	}

	// Get request ID
	reqId := middleware.GetReqID(r.Context())

	// Marshall response
	var response = map[string]interface{}{
		"message": m,
		"request": reqId,
	}
	if e != nil {
		response["error"] = e.Error()
	}
	json, err := json.Marshal(response)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Error marshalling error response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	w.Write(json)
}
