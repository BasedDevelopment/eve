package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	req := new(util.LoginRequest)

	if err := util.ParseRequest(r, req); err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "Failed to parse login request")
		return
	}

	// New profile instance
	profile := controllers.Profile{Email: req.Email}
	profile, err := profile.Get(ctx)

	if err != nil {
		util.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(req.Password)); err != nil {
		util.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Issue token
	userToken, err := sessions.NewSession(ctx, profile)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Send token to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": userToken.String(),
	})
}
