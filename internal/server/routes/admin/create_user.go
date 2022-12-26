package admin

import (
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request
	req := new(util.CreateRequest)

	if err := util.ParseRequest(r, req); err != nil {
		util.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	// New profile instance
	profile := controllers.Profile{
		Email:    req.Email,
		Name:     req.Name,
		Disabled: req.Disabled,
		IsAdmin:  req.IsAdmin,
		Remarks:  req.Remarks,
	}

	if profile, _ := profile.Get(ctx); profile.Name != "" {
		util.WriteError(w, r, nil, http.StatusBadRequest, "User already exists")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	profile.Password = string(hash)
	uuid, err := profile.New(ctx)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to create user")
		return
	}

	resp := (map[string]interface{}{
		"uuid": uuid,
	})

	util.WriteResponse(resp, w, http.StatusCreated)
}
