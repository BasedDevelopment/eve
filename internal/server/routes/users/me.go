package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p := new(controllers.Profile)
	id := ctx.Value("id").(string)
	p.ID = uuid.MustParse(id)

	profile, err := p.Get(ctx)
	if err != nil {
		if errors.Is(err, controllers.QueryErr) {
			log.Error().Err(err).Msg("query error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		if errors.Is(err, controllers.CollectErr) {
			log.Error().Err(err).Msg("collect error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
	}
	type response struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		LastLogin time.Time `json:"lastLogin"`
		Created   time.Time `json:"created"`
		Updated   time.Time `json:"updated"`
	}

	outJson, err := json.Marshal(response{
		Name:      profile.Name,
		Email:     profile.Email,
		LastLogin: profile.LastLogin,
		Created:   profile.Created,
		Updated:   profile.Updated,
	})
	if err != nil {
		log.Error().Err(err).Msg("marshal error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(outJson)
}
