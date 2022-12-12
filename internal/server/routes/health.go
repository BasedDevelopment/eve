package routes

import (
	"net/http"

	"github.com/ericzty/eve/internal/util"
)

func Health(w http.ResponseWriter, r *http.Request) {
	t := new(util.LoginRequest)

	if err := util.ParseRequest(r, t); err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
		return
	}

	util.WriteResponse(t, w, http.StatusOK)
}
