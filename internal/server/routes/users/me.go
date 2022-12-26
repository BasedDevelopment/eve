package users

import (
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
	"github.com/google/uuid"
)

func GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile := controllers.Profile{ID: ctx.Value("owner").(uuid.UUID)}
	profile, err := profile.Get(ctx)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = util.WriteResponse(util.UserResponse{
		Name:      profile.Name,
		Email:     profile.Email,
		LastLogin: profile.LastLogin,
		Created:   profile.Created,
		Updated:   profile.Updated,
	}, w, http.StatusOK)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
