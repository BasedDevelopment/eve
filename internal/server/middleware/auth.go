package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/tokens"
	"github.com/rs/zerolog/log"
)

// Auth forces a user to be authenticated before continuing to the route
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing Auth-Token header"))

			return
		}

		splitToken := strings.Split(authorizationHeader, "Bearer ")
		reqToken := tokens.Parse(splitToken[1])

		isValidSession := sessions.ValidateSession(ctx, reqToken)

		if !isValidSession {
		}

		isAdmin, err := controllers.IsAdmin(ctx, id)

		if err != nil {
			log.Error().Err(err).Msg("Error fetching isAdmin from db")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))

			return
		}

		ctx = context.WithValue(ctx, "id", id)
		ctx = context.WithValue(ctx, "isAdmin", isAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAdmin := ctx.Value("isAdmin").(bool)

		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))

			return
		}

		next.ServeHTTP(w, r)
	})
}
