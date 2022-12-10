package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/sessions"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	// Check request body
	if loginRequest.Email == "" || loginRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing required fields"))
		return
	}

	// New profile instance
	profile := controllers.Profile{Email: loginRequest.Email}
	profile, err := profile.Get(ctx)

	// Validate password
	hash, err := profile.GetHash(ctx)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(loginRequest.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	// Issue token
	userToken, err := sessions.NewSession(ctx, profile)

	if err != nil {
		log.Error().Err(err).Msg("Failed to issue token")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	// Send token to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": userToken.String(),
	})
}
