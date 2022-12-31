package routes

import (
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
)

func GetVMs(w http.ResponseWriter, r *http.Request) {
	// Get hv ID from request
	hvid := strings.Split(r.URL.Path, "/")[3]
	hv := controllers.Cloud.HVs[hvid]
	var vms []*controllers.VM
	for _, vm := range hv.VMs {
		vms = append(vms, vm)
	}

	// Send response
	if err := util.WriteResponse(vms, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}

}
