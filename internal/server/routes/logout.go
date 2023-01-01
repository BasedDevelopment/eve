package routes

import (
	"net/http"
	"strings"

	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/tokens"
	"github.com/BasedDevelopment/eve/internal/util"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	reqToken, err := tokens.Parse(header[1])
	if err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "Invalid token")
		return
	}

	sessions.Delete(ctx, reqToken)

	if err := sessions.Delete(ctx, reqToken); err != nil {
		util.WriteError(w, r, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.WriteResponse(map[string]interface{}{
		"message": "logout success",
	}, w, http.StatusOK)
}
