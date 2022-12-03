package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//make create user struct
	var createRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Disabled bool   `json:"disabled"`
		IsAdmin  bool   `json:"is_admin"`
		Remarks  string `json:"remarks"`
	}

	// Decode request body see what user wants
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	// Check request body see if any are empty
	if createRequest.Email == "" || createRequest.Password == "" || createRequest.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	// New profile instance
	p := new(controllers.Profile)
	p.Email = createRequest.Email
	p.Name = createRequest.Name
	p.Disabled = createRequest.Disabled
	p.IsAdmin = createRequest.IsAdmin
	if createRequest.Remarks != "" {
		p.Remarks.String = createRequest.Remarks
		p.Remarks.Status = 2
	}

	existingUserHash, err := p.GetHash(ctx)

	if err != nil {
		if err.Error() != "no rows in result set" {
			log.Error().Err(err).Msg("Failed to get hash")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
	}

	if existingUserHash != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User already exists"))
		return
	}

	// Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(createRequest.Password), 10)
	hash := string(bytes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	p.Password = hash

	uuid, err := p.New(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(uuid))
		return
	}
}
