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
	reqToken := tokens.Parse(header[1])

	sessions.Delete(ctx, reqToken)

	err := sessions.Delete(ctx, reqToken)

	if err != nil {
		util.WriteError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout Success"))

	return
}
