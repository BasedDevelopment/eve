package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/util"
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
		util.WriteError(err, w, http.StatusBadRequest)
		return
	}

	// Check request body
	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.WriteError(errors.New("Missing required fields"), w, http.StatusBadRequest)
		return
	}

	// New profile instance
	profile := controllers.Profile{Email: loginRequest.Email}
	profile, err := profile.Get(ctx)

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(loginRequest.Password)); err != nil {
		util.WriteError(errors.New("Unauthorized"), w, http.StatusUnauthorized)
		return
	}

	// Issue token
	userToken, err := sessions.NewSession(ctx, profile)

	if err != nil {
		util.WriteError(errors.New("Failed to issue token"), w, http.StatusInternalServerError)
		return
	}

	// Send token to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": userToken.String(),
	})
}
