package admin

import (
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
)

func GetHVs(w http.ResponseWriter, r *http.Request) {
	cloud := controllers.Cloud
	var out []*controllers.HV
	for _, hv := range cloud.HVs {
		hv.VMs = make(map[string]*controllers.VM)
		out = append(out, hv)
	}

	// Send response
	if err := util.WriteResponse(out, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetHV(w http.ResponseWriter, r *http.Request) {
	// Get hv ID from request
	parts := strings.Split(r.URL.Path, "/")
	hvid := parts[3]

	hv := controllers.Cloud.HVs[hvid]

	// Send response
	if err := util.WriteResponse(hv, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}
