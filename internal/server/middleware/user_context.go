package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/util"
	"github.com/rs/zerolog/log"
)

// UserContext fetches the owner of the request from the current session and
// and appends it to the request context. Requires Auth, required MustBeAdmin.
func UserContext(next http.Handler) http.Handler {
	// This function doesn't check whether a user is authenticated
	// and as such should only be used after Auth has been called.

	// It is required for the MustBeAdmin middleware though, since
	// that middleware uses the profile in the request context.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get session
		requestToken := getToken(w, r) // function from auth middleware; ges token from authorization header
		session, err := sessions.GetSession(ctx, requestToken)

		if err != nil {
			log.Error().Err(err).Msg("Error while setting user context")

			util.WriteError(errors.New("Internal Server Error"), w, http.StatusInternalServerError)

			return
		}

		// Add the owner of the session to r.Context
		ctx = context.WithValue(ctx, "owner", session.Owner)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
