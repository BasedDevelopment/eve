package util

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	LastLogin time.Time `json:"lastLogin"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

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

func WriteError(w http.ResponseWriter, r *http.Request, e error, s int, m string) {
	if e != nil {
		log.Error().
			Err(e).
			Str("message", m).
			Msg("Request Error")
	}

	// Get request ID
	reqId := middleware.GetReqID(r.Context())

	// Marshall response
	json, err := json.Marshal(map[string]interface{}{
		"message":   m,
		"error":     e,
		"requestID": reqId,
	})

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
