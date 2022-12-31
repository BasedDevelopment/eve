package admin

import (
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
	"github.com/google/uuid"
)

func GetHVs(w http.ResponseWriter, r *http.Request) {
	cloud := controllers.Cloud
	var vms []*controllers.HV
	for _, hv := range cloud.HVs {
		hv.VMs = make(map[uuid.UUID]*controllers.VM)
		vms = append(vms, hv)
	}

	// Send response
	if err := util.WriteResponse(vms, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetHV(w http.ResponseWriter, r *http.Request) {
	// Get hv ID from request
	hvid := strings.Split(r.URL.Path, "/")[3]

	hv := controllers.Cloud.HVs[hvid]

	// Send response
	if err := util.WriteResponse(hv, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}
