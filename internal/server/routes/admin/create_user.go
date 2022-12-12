package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Disabled bool   `json:"disabled"`
		IsAdmin  bool   `json:"is_admin"`
		Remarks  string `json:"remarks"`
	}

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&createRequest); err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
		return
	}

	// Check request body see if any are empty
	if createRequest.Email == "" || createRequest.Password == "" || createRequest.Name == "" {
		util.WriteError(errors.New("Bad Request"), w, http.StatusBadRequest)
		return
	}

	// New profile instance
	profile := controllers.Profile{
		Email:     createRequest.Email,
		Name:      createRequest.Name,
		Disabled:  createRequest.Disabled,
		IsAdmin:   createRequest.IsAdmin,
		Remarks:   createRequest.Remarks,
		LastLogin: time.Now(),
	}

	// Check if user exists

	// if err != nil {
	// 	if err.Error() != "no rows in result set" {
	// 		log.Error().Err(err).Msg("Failed to get hash")
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte("Internal server error"))

	// 		return
	// 	}
	// }

	// if existingUserHash != "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("User already exists"))

	// 	return
	// }

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(createRequest.Password), 10)

	if err != nil {
		util.WriteError(err, w, http.StatusInternalServerError)
		return
	}

	profile.Password = string(hash)
	uuid, err := profile.New(ctx)

	if err != nil {
		util.WriteError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(uuid))
	return
}
