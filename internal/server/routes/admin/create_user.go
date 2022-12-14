package admin

import (
	"net/http"
	"time"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request
	req := new(util.CreateRequest)

	if err := util.ParseRequest(r, req); err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
	}

	// New profile instance
	profile := controllers.Profile{
		Email:     req.Email,
		Name:      req.Name,
		Disabled:  req.Disabled,
		IsAdmin:   req.IsAdmin,
		Remarks:   req.Remarks,
		LastLogin: time.Now(),
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

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

	util.WriteResponse(map[string]interface{}{
		"id": uuid,
	}, w, http.StatusCreated)

	return
}
