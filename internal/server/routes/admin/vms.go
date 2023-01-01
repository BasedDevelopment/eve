package admin

import (
	"net/http"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/go-chi/chi/v5"
)

func GetVMs(w http.ResponseWriter, r *http.Request) {
	// Get hv ID from request
	hvid := chi.URLParam(r, "id")
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
