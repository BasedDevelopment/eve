package routes

import (
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/sessions"
	"github.com/ericzty/eve/internal/tokens"
	"github.com/ericzty/eve/internal/util"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	reqToken, err := tokens.Parse(header[1])
	if err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
		return
	}

	sessions.Delete(ctx, reqToken)

	if err := sessions.Delete(ctx, reqToken); err != nil {
		util.WriteError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout Success"))

	return
}
