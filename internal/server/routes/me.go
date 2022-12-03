package routes

import (
	"errors"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p := new(controllers.Profile)
	id := ctx.Value("id").(string)
	p.ID = uuid.MustParse(id)

	if err := p.Get(ctx); err != nil {
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
}
