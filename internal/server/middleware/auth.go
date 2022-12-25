package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/tokens"
	"github.com/ericzty/eve/internal/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func getToken(w http.ResponseWriter, r *http.Request) (token tokens.Token) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		util.WriteError(errors.New("Missing Authorization header"), w, http.StatusBadRequest)

		return tokens.Token{}
	}

	splitHeader := strings.Split(authorizationHeader, "Bearer ")
	token, err := tokens.Parse(splitHeader[1])
	if err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
		return
	}
	return
}

// Auth forces a user to be authenticated before continuing to the route
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestToken := getToken(w, r)

		if !sessions.ValidateSession(ctx, requestToken) {
			util.WriteError(errors.New("Unauthorized"), w, http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MustBeAdmin verifies is an **already authenticated** user is an admin
func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		profile := controllers.Profile{ID: ctx.Value("owner").(uuid.UUID)}
		profile, err := profile.Get(ctx)

		if err != nil {
			log.Error().Err(err).Msg("User Fetch")
			util.WriteError(errors.New("Internal Server Error"), w, http.StatusInternalServerError)

			return
		}

		if !profile.IsAdmin {
			util.WriteError(errors.New("Unauthorized"), w, http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}
