package routes

import (
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/rs/zerolog/log"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	if err := controllers.Logout(ctx, reqToken); err != nil {
		log.Error().Err(err).Msg("Logout")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logout Success"))
		return
	}
}
