package util

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type UserResponse struct {
	Name      string
	Email     string
	LastLogin time.Time
	Created   time.Time
	Updated   time.Time
}

func WriteResponse[R UserResponse | map[string]interface{} | []map[string]interface{}](r R, w http.ResponseWriter, status ...int) error {
	json, err := json.Marshal(r)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	if status != nil {
		w.WriteHeader(status[0])
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return nil
}

func WriteError(e error, w http.ResponseWriter, s int) {
	log.Debug().Err(e).Msg("Request Error")

	json, err := json.Marshal(map[string]interface{}{
		"error": e.Error(),
	})

	log.Debug().Err(err).Msg("Marshal Error")

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	w.WriteHeader(s)
}
