package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/tokens"
	"github.com/gofrs/uuid"
)

func getToken(w http.ResponseWriter, r *http.Request) tokens.Token {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing Authorization header"))

		return tokens.Token{}
	}

	splitHeader := strings.Split(authorizationHeader, "Bearer ")
	return tokens.Parse(splitHeader[1])
}

// Auth forces a user to be authenticated before continuing to the route
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestToken := getToken(w, r)

		if !sessions.ValidateSession(ctx, requestToken) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))

			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MustBeAdmin verifies is an **already authenticated** user is an admin
func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		owner := ctx.Value("owner").(uuid.UUID)

		fmt.Println(owner)

		next.ServeHTTP(w, r)
	})
}
