package middlewares

import (
	"net/http"
	"strings"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/tokens"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
)

func getToken(w http.ResponseWriter, r *http.Request) (token tokens.Token) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		util.WriteError(w, r, nil, http.StatusBadRequest, "Missing Authorization header")
		return tokens.Token{}
	}

	splitHeader := strings.Split(authorizationHeader, "Bearer ")
	token, err := tokens.Parse(splitHeader[1])
	if err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "Invalid token")
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
			util.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")

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
			util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")

			return
		}

		if !profile.IsAdmin {
			util.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")

			return
		}

		next.ServeHTTP(w, r)
	})
}
