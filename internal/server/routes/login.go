package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	req := new(util.LoginRequest)

	if err := util.ParseRequest(r, req); err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
		return
	}

	// New profile instance
	profile := controllers.Profile{Email: loginRequest.Email}
	profile, err := profile.Get(ctx)

	if err != nil {
		log.Debug().Err(err).Msg("Profile Fetch Error")
		util.WriteError(err, w, http.StatusNotFound)
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(loginRequest.Password)); err != nil {
		log.Debug().Err(err).Msg("Comparison Error")
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
