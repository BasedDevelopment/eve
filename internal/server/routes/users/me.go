package users

import (
	"encoding/json"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	profile := controllers.Profile{ID: ctx.Value("owner").(uuid.UUID)}
	profile, err := profile.Get(ctx)

	if err != nil {
		log.Error().Err(err).Msg("User Fetch")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))

		return
	}

	outJson, err := json.Marshal(map[string]interface{}{
		"name":       profile.Name,
		"email":      profile.Email,
		"last_login": profile.LastLogin,
		"created":    profile.Created,
		"updated":    profile.Updated,
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
